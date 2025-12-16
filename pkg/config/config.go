// Package config provides the configuration for the application.
package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kelseyhightower/envconfig"
)

// Config represents the configuration for the application
type Config struct {
	Timeout    time.Duration  `envconfig:"TIMEOUT" default:"180s"`
	LogLevel   string         `envconfig:"LOG_LEVEL" default:"info"`
	AWSConfig  aws.Config     // Loaded using AWS SDK, not from env
	Postgres   PostgresConfig `envconfig:"POSTGRES"`
	SQLitePath string         `envconfig:"SQLITE_PATH" default:"./data/assistant.db"`

	// AI configuration (organized by provider)
	AI AIConfig `envconfig:"AI"`

	// Records configuration
	Sources SourcesConfig `envconfig:"SOURCES"`
}

// OllamaConfig represents the configuration for local AI services
type OllamaConfig struct {
	URL   string `envconfig:"URL" default:"http://localhost:11434"`
	Model string `envconfig:"MODEL" default:"codellama:7b-instruct"`
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
	Ollama  OllamaConfig    `envconfig:"OLLAMA"`
	Bedrock BedrockAIConfig `envconfig:"BEDROCK"`
}

// RDSPostgres represents the configuration for AWS RDS Postgres
type RDSPostgres struct {
	CredentialsSecretARN  string `envconfig:"CREDENTIALS_SECRET_ARN"`
	SchemaEnsureLambdaARN string `envconfig:"RDS_POSTGRES_SCHEMA_ENSURE_LAMBDA_ARN"`
	InstanceARN           string `envconfig:"INSTANCE_ARN"`
	DatabaseName          string `envconfig:"DATABASE_NAME" default:"assistant_db"`
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

// SourcesConfig represents configuration for data sources
type SourcesConfig struct {
	StoragePath string            `envconfig:"STORAGE_PATH" default:"./data/records"`
	Local       LocalSourceConfig `envconfig:"LOCAL"`
}

// LocalSourceConfig represents configuration for local file source
type LocalSourceConfig struct {
	Enabled  bool   `envconfig:"ENABLED" default:"true"`
	BasePath string `envconfig:"BASE_PATH" default:"./testdata"`
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
	var cfg Config

	// Load env vars
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Setup structured logging as early as possible
	setupLogger(cfg.LogLevel)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return cfg, fmt.Errorf("failed to load AWS configuration: %w", err)
	}
	cfg.AWSConfig = awsCfg

	return cfg, nil
}
