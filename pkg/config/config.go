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
	Timeout    time.Duration `env:"TIMEOUT" envDefault:"180s"`
	LogLevel   string        `env:"LOG_LEVEL" envDefault:"info"`
	AWSConfig  aws.Config    // Loaded using AWS SDK, not from env
	SQLitePath string        `env:"SQLITE_PATH" envDefault:"./data/assistant.db"`

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

// AIConfig represents the overall AI configuration with provider-specific settings
type AIConfig struct {
	// Provider selection (can be overridden per request)
	DefaultProvider string `env:"DEFAULT_PROVIDER" envDefault:"bedrock"`

	// Provider-specific configurations
	Ollama OllamaConfig `envPrefix:"OLLAMA_"`
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
