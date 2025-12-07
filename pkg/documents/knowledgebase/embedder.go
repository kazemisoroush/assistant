package knowledgebase

import "context"

// Embedder generates vector embeddings from text
type Embedder interface {
	// Embed generates embeddings for text
	// Returns a vector of floats representing the semantic meaning
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch generates embeddings for multiple texts
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Dimensions returns the dimension of the embedding vectors
	Dimensions() int
}

// EmbedderConfig represents configuration for an embedder
type EmbedderConfig struct {
	Provider string // "ollama", "bedrock", "openai", etc.
	Model    string // Model name
	APIKey   string // API key if required
	Endpoint string // Custom endpoint if required
}
