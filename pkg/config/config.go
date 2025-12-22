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
	"github.com/caarlos0/env/v11"
)

// Config represents the configuration for the application
type Config struct {
	Timeout    time.Duration  `env:"TIMEOUT" envDefault:"180s"`
	LogLevel   string         `env:"LOG_LEVEL" envDefault:"info"`
	AWSConfig  aws.Config     // Loaded using AWS SDK, not from env
	Postgres   PostgresConfig `envPrefix:"POSTGRES_"`
	SQLitePath string         `env:"SQLITE_PATH" envDefault:"./data/assistant.db"`

	// AI configuration (organized by provider)
	AI AIConfig `envPrefix:"AI_"`

	// Records configuration
	Sources SourcesConfig `envPrefix:"SOURCES_"`
}

// OllamaConfig represents the configuration for local AI services
type OllamaConfig struct {
	URL   string `env:"URL" envDefault:"http://localhost:11434"`
	Model string `env:"MODEL" envDefault:"codellama:7b-instruct"`
}

// BedrockAIConfig represents the configuration for AWS Bedrock AI services
type BedrockAIConfig struct {
	Region                      string      `env:"REGION" envDefault:"us-east-1"`
	KnowledgeBaseServiceRoleARN string      `env:"KNOWLEDGE_BASE_SERVICE_ROLE_ARN"`
	AgentServiceRoleARN         string      `env:"AGENT_SERVICE_ROLE_ARN"`
	FoundationModel             string      `env:"FOUNDATION_MODEL" envDefault:"amazon.titan-tg1-large"`
	S3BucketName                string      `env:"S3_BUCKET_NAME"`
	RDSPostgres                 RDSPostgres `envPrefix:"RDS_POSTGRES_"`
}

// AIConfig represents the overall AI configuration with provider-specific settings
type AIConfig struct {
	// Provider selection (can be overridden per request)
	DefaultProvider string `env:"DEFAULT_PROVIDER" envDefault:"bedrock"`

	// Provider-specific configurations
	Ollama  OllamaConfig    `envPrefix:"OLLAMA_"`
	Bedrock BedrockAIConfig `envPrefix:"BEDROCK_"`
}

// RDSPostgres represents the configuration for AWS RDS Postgres
type RDSPostgres struct {
	CredentialsSecretARN  string `env:"CREDENTIALS_SECRET_ARN"`
	SchemaEnsureLambdaARN string `env:"SCHEMA_ENSURE_LAMBDA_ARN"`
	InstanceARN           string `env:"INSTANCE_ARN"`
	DatabaseName          string `env:"DATABASE_NAME" envDefault:"assistant_db"`
}

// PostgresConfig represents the configuration for PostgreSQL connection
type PostgresConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	Database string `env:"DATABASE" envDefault:"assistant_db"`
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

// SourcesConfig represents configuration for data sources
type SourcesConfig struct {
	StoragePath string            `env:"STORAGE_PATH" envDefault:"./data/records"`
	Local       LocalSourceConfig `envPrefix:"LOCAL_"`
}

// LocalSourceConfig represents configuration for local file source
type LocalSourceConfig struct {
	Enabled  bool   `env:"ENABLED" envDefault:"true"`
	BasePath string `env:"BASE_PATH" envDefault:"./testdata"`
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
	if err := env.Parse(&cfg); err != nil {
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
