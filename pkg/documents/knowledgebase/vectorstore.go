package knowledgebase

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/documents"
)

// VectorStore defines operations for vector-based document search
// This is an interface for future implementation with Chroma, Pinecone, or AWS Bedrock
type VectorStore interface {
	// Index adds document embeddings to the vector store
	Index(ctx context.Context, doc *documents.Document) error

	// Search performs semantic similarity search
	Search(ctx context.Context, query string, limit int) ([]documents.SearchResult, error)

	// Delete removes document from vector store
	Delete(ctx context.Context, docID string) error

	// Close closes the vector store connection
	Close() error
}

// TODO: Implement concrete implementations:
// - OllamaEmbedder: Use local Ollama for embeddings
// - BedrockEmbedder: Use AWS Bedrock for embeddings
// - ChromaVectorStore: Use Chroma for vector storage
// - LocalVectorStore: Simple in-memory vector store for development
