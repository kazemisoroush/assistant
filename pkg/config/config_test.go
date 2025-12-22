package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_Success tests loading configuration with valid environment variables
func TestLoadConfig_Success(t *testing.T) {
	// Setup environment variables
	envVars := map[string]string{
		"TIMEOUT":                     "120s",
		"LOG_LEVEL":                   "debug",
		"SQLITE_PATH":                 "/tmp/test.db",
		"AI_DEFAULT_PROVIDER":         "ollama",
		"AI_OLLAMA_URL":               "http://localhost:11434",
		"AI_OLLAMA_MODEL":             "llama2",
		"AI_BEDROCK_REGION":           "us-west-2",
		"AI_BEDROCK_FOUNDATION_MODEL": "anthropic.claude-v2",
		"AI_BEDROCK_S3_BUCKET_NAME":   "test-bucket",
		"POSTGRES_HOST":               "db.example.com",
		"POSTGRES_PORT":               "5433",
		"POSTGRES_DATABASE":           "testdb",
		"POSTGRES_USERNAME":           "testuser",
		"POSTGRES_PASSWORD":           "testpass",
		"POSTGRES_SSL_MODE":           "require",
		"SOURCES_STORAGE_PATH":        "/data/test",
		"SOURCES_LOCAL_ENABLED":       "true",
		"SOURCES_LOCAL_BASE_PATH":     "/tmp/testdata",
	}

	// Set environment variables
	for key, value := range envVars {
		err := os.Setenv(key, value)
		require.NoError(t, err, "Failed to set env var %s", key)
	}

	// Cleanup after test
	defer func() {
		for key := range envVars {
			_ = os.Unsetenv(key)
		}
	}()

	// Load configuration
	cfg, err := LoadConfig()
	require.NoError(t, err, "LoadConfig() should not fail")

	// Validate configuration values
	assert.Equal(t, 120*time.Second, cfg.Timeout, "Timeout should be 120s")
	assert.Equal(t, "debug", cfg.LogLevel, "LogLevel should be 'debug'")
	assert.Equal(t, "/tmp/test.db", cfg.SQLitePath, "SQLitePath should be '/tmp/test.db'")

	// AI configuration
	assert.Equal(t, "ollama", cfg.AI.DefaultProvider, "AI.DefaultProvider should be 'ollama'")
	assert.Equal(t, "http://localhost:11434", cfg.AI.Ollama.URL, "AI.Ollama.URL should be 'http://localhost:11434'")
	assert.Equal(t, "llama2", cfg.AI.Ollama.Model, "AI.Ollama.Model should be 'llama2'")
	assert.Equal(t, "us-west-2", cfg.AI.Bedrock.Region, "AI.Bedrock.Region should be 'us-west-2'")
	assert.Equal(t, "anthropic.claude-v2", cfg.AI.Bedrock.FoundationModel, "AI.Bedrock.FoundationModel should be 'anthropic.claude-v2'")

	// Postgres configuration
	assert.Equal(t, "db.example.com", cfg.Postgres.Host, "Postgres.Host should be 'db.example.com'")
	assert.Equal(t, 5433, cfg.Postgres.Port, "Postgres.Port should be 5433")
	assert.Equal(t, "testdb", cfg.Postgres.Database, "Postgres.Database should be 'testdb'")
	assert.Equal(t, "testuser", cfg.Postgres.Username, "Postgres.Username should be 'testuser'")
	assert.Equal(t, "testpass", cfg.Postgres.Password, "Postgres.Password should be 'testpass'")
	assert.Equal(t, "require", cfg.Postgres.SSLMode, "Postgres.SSLMode should be 'require'")

	// Sources configuration
	assert.Equal(t, "/data/test", cfg.Sources.StoragePath, "Sources.StoragePath should be '/data/test'")
	assert.True(t, cfg.Sources.Local.Enabled, "Sources.Local.Enabled should be true")
	assert.Equal(t, "/tmp/testdata", cfg.Sources.Local.BasePath, "Sources.Local.BasePath should be '/tmp/testdata'")

	// Verify AWS config was loaded (should not be nil/zero value)
	if cfg.AWSConfig.Region == "" {
		t.Log("Warning: AWS config region is empty (may be expected in test environment)")
	}
}

// TestLoadConfig_DefaultValues tests that default values are applied when env vars are not set
func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear all relevant environment variables to ensure defaults are used
	envVarsToClear := []string{
		"TIMEOUT",
		"LOG_LEVEL",
		"SQLITE_PATH",
		"AI_DEFAULT_PROVIDER",
		"AI_OLLAMA_URL",
		"AI_OLLAMA_MODEL",
		"AI_BEDROCK_REGION",
		"AI_BEDROCK_FOUNDATION_MODEL",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DATABASE",
		"POSTGRES_USERNAME",
		"POSTGRES_PASSWORD",
		"POSTGRES_SSL_MODE",
		"SOURCES_STORAGE_PATH",
		"SOURCES_LOCAL_ENABLED",
		"SOURCES_LOCAL_BASE_PATH",
	}

	for _, key := range envVarsToClear {
		_ = os.Unsetenv(key)
	}

	// Load configuration
	cfg, err := LoadConfig()
	require.NoError(t, err, "LoadConfig() should not fail")

	// Validate default values
	assert.Equal(t, 180*time.Second, cfg.Timeout, "Default Timeout should be 180s")
	assert.Equal(t, "info", cfg.LogLevel, "Default LogLevel should be 'info'")
	assert.Equal(t, "./data/assistant.db", cfg.SQLitePath, "Default SQLitePath should be './data/assistant.db'")

	// AI configuration defaults
	assert.Equal(t, "bedrock", cfg.AI.DefaultProvider, "Default AI.DefaultProvider should be 'bedrock'")
	assert.Equal(t, "http://localhost:11434", cfg.AI.Ollama.URL, "Default AI.Ollama.URL should be 'http://localhost:11434'")
	assert.Equal(t, "codellama:7b-instruct", cfg.AI.Ollama.Model, "Default AI.Ollama.Model should be 'codellama:7b-instruct'")
	assert.Equal(t, "us-east-1", cfg.AI.Bedrock.Region, "Default AI.Bedrock.Region should be 'us-east-1'")
	assert.Equal(t, "amazon.titan-tg1-large", cfg.AI.Bedrock.FoundationModel, "Default AI.Bedrock.FoundationModel should be 'amazon.titan-tg1-large'")

	// Postgres configuration defaults
	assert.Equal(t, "localhost", cfg.Postgres.Host, "Default Postgres.Host should be 'localhost'")
	assert.Equal(t, 5432, cfg.Postgres.Port, "Default Postgres.Port should be 5432")
	assert.Equal(t, "assistant_db", cfg.Postgres.Database, "Default Postgres.Database should be 'assistant_db'")
	assert.Equal(t, "postgres", cfg.Postgres.Username, "Default Postgres.Username should be 'postgres'")
	assert.Equal(t, "disable", cfg.Postgres.SSLMode, "Default Postgres.SSLMode should be 'disable'")

	// Sources configuration defaults
	assert.Equal(t, "./data/records", cfg.Sources.StoragePath, "Default Sources.StoragePath should be './data/records'")
	assert.True(t, cfg.Sources.Local.Enabled, "Default Sources.Local.Enabled should be true")
	assert.Equal(t, "./testdata", cfg.Sources.Local.BasePath, "Default Sources.Local.BasePath should be './testdata'")
}
