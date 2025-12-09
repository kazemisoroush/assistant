// Package config provides the configuration for the application.
package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration for the application
type Config struct {
	Git            GitConfig      `envconfig:"GIT"`
	TimeoutSeconds int            `envconfig:"TIMEOUT_SECONDS" default:"180"`
	LogLevel       string         `envconfig:"LOG_LEVEL" default:"info"`
	AWSConfig      aws.Config     // Loaded using AWS SDK, not from env
	Cognito        CognitoConfig  `envconfig:"COGNITO"`
	Postgres       PostgresConfig `envconfig:"POSTGRES"`

	// AI configuration (organized by provider)
	AI AIConfig `envconfig:"AI"`

	// Record storage configuration
	Records RecordsConfig `envconfig:"RECORDS"`
}

// LocalAIConfig represents the configuration for local AI services
type LocalAIConfig struct {
	Enabled        bool   `envconfig:"ENABLED" default:"false"`
	OllamaURL      string `envconfig:"OLLAMA_URL" default:"http://localhost:11434"`
	Model          string `envconfig:"MODEL" default:"codellama:7b-instruct"`
	ChromaURL      string `envconfig:"CHROMA_URL" default:"http://localhost:8000"`
	EmbeddingModel string `envconfig:"EMBEDDING_MODEL" default:"all-MiniLM-L6-v2"`
}

// BedrockAIConfig represents the configuration for AWS Bedrock AI services
type BedrockAIConfig struct {
	Region                      string      `envconfig:"REGION" default:"us-east-1"`
	KnowledgeBaseServiceRoleARN string      `envconfig:"KNOWLEDGE_BASE_SERVICE_ROLE_ARN"`
	AgentServiceRoleARN         string      `envconfig:"AGENT_SERVICE_ROLE_ARN"`
	FoundationModel             string      `envconfig:"FOUNDATION_MODEL" default:"amazon.titan-tg1-large"`
	S3BucketName                string      `envconfig:"S3_BUCKET_NAME"`
	RDSPostgres                 RDSPostgres `envconfig:"RDS_POSTGRES"`
}

// AIConfig represents the overall AI configuration with provider-specific settings
type AIConfig struct {
	// Provider selection (can be overridden per request)
	DefaultProvider string `envconfig:"DEFAULT_PROVIDER" default:"bedrock"`

	// Provider-specific configurations
	Local   LocalAIConfig   `envconfig:"LOCAL"`
	Bedrock BedrockAIConfig `envconfig:"BEDROCK"`
}

// RDSPostgres represents the configuration for AWS RDS Postgres
type RDSPostgres struct {
	CredentialsSecretARN  string `envconfig:"CREDENTIALS_SECRET_ARN"`
	SchemaEnsureLambdaARN string `envconfig:"RDS_POSTGRES_SCHEMA_ENSURE_LAMBDA_ARN"`
	InstanceARN           string `envconfig:"INSTANCE_ARN"`
	DatabaseName          string `envconfig:"DATABASE_NAME" default:"assistant_db"`
}

// CognitoConfig represents the configuration for AWS Cognito authentication
type CognitoConfig struct {
	UserPoolID string `envconfig:"USER_POOL_ID" required:"true"`
	ClientID   string `envconfig:"CLIENT_ID" required:"true"`
	Region     string `envconfig:"REGION" default:"us-east-1"`
}

// PostgresConfig represents the configuration for PostgreSQL connection
type PostgresConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	Database string `envconfig:"DATABASE" default:"assistant_db"`
	Username string `envconfig:"USERNAME" default:"postgres"`
	Password string `envconfig:"PASSWORD"`
	SSLMode  string `envconfig:"SSL_MODE" default:"disable"`
}

// DatabaseSecret represents the structure of the secret stored in AWS Secrets Manager
type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"dbname"`
}

// GitConfig represents the Git configuration
type GitConfig struct {
	CodebaseURL string `envconfig:"CODEBASE_URL"` // Per-request, not required at startup
	Token       string `envconfig:"TOKEN" required:"true"`
	Author      string `envconfig:"AUTHOR" default:"AssistantBot"`
	Email       string `envconfig:"EMAIL" default:"bot@example.com"`
}

// RecordsConfig represents configuration for record storage and processing
type RecordsConfig struct {
	StoragePath string `envconfig:"STORAGE_PATH" default:"./data/records"`
}

// VectorStoreConfig represents configuration for a vector store
type VectorStoreConfig struct {
	Provider string // "chroma", "pinecone", "bedrock", "local", etc.
	Endpoint string // Connection endpoint
	APIKey   string // API key if required
	Index    string // Index/collection name
}

// validateRepositoryURL ensures the RepoURL matches the expected GitHub URL pattern
func validateRepositoryURL(url string) error {
	// Regex for GitHub repo URL (HTTPS and SSH formats)
	gitHubURLRegex := `^(https:\/\/github\.com\/[\w-]+\/[\w.-]+(\.git)?|git@github\.com:[\w-]+\/[\w.-]+(\.git)?)$`
	matched, err := regexp.MatchString(gitHubURLRegex, url)
	if err != nil {
		return fmt.Errorf("failed to validate GitHub URL: %w", err)
	}
	if !matched {
		return errors.New("invalid GitHub repository URL format")
	}
	return nil
}

// setupLogger configures slog with JSON output and the specified log level
func setupLogger(level string) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create JSON handler with specified log level
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	// Set the default logger
	slog.SetDefault(slog.New(handler))
}

// LoadConfig loads and validates configuration from environment variables and AWS
func LoadConfig() (Config, error) {
	return LoadConfigWithDependencies()
}

// LoadConfigWithDependencies loads configuration with optional dependency injection for testing
func LoadConfigWithDependencies() (Config, error) {
	var cfg Config

	// Load env vars
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Setup structured logging as early as possible
	setupLogger(cfg.LogLevel)

	// Validate RepoURL if provided (it's optional at startup, provided per-request)
	if cfg.Git.CodebaseURL != "" {
		if err := validateRepositoryURL(cfg.Git.CodebaseURL); err != nil {
			return cfg, fmt.Errorf("invalid GitHub repository URL: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSeconds)*time.Second)
	defer cancel()

	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return cfg, fmt.Errorf("failed to load AWS configuration: %w", err)
	}
	cfg.AWSConfig = awsCfg

	return cfg, nil
}
