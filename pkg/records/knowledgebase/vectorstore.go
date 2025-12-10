package knowledgebase

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// VectorStore defines operations for vector-based record search
// This is an interface for future implementation with Chroma, Pinecone, or AWS Bedrock
//
//go:generate mockgen -destination=./mocks/mock_vectorstore.go -mock_names=VectorStore=MockVectorStore -package=mocks . VectorStore
type VectorStore interface {
	// Index adds record embeddings to the vector store
	Index(ctx context.Context, rec *records.Record) error

	// Search performs semantic similarity search
	Search(ctx context.Context, query string, limit int) ([]records.SearchResult, error)

	// Delete removes record from vector store
	Delete(ctx context.Context, recID string) error

	// Close closes the vector store connection
	Close() error
}

// TODO: Implement concrete implementations:
// - OllamaEmbedder: Use local Ollama for embeddings
// - BedrockEmbedder: Use AWS Bedrock for embeddings
// - ChromaVectorStore: Use Chroma for vector storage
// - LocalVectorStore: Simple in-memory vector store for development
