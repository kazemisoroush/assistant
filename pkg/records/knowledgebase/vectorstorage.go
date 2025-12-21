package knowledgebase

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// VectorStorage defines operations for vector-based record search
// This is an interface for future implementation with Chroma, Pinecone, or AWS Bedrock
//
//go:generate mockgen -destination=./mocks/mock_vectorstorage.go -mock_names=VectorStorage=MockVectorStorage -package=mocks . VectorStorage
type VectorStorage interface {
	// Index adds record embeddings to the vector store
	Index(ctx context.Context, rec records.Record) error

	// Search performs semantic similarity search
	Search(ctx context.Context, prompt string) ([]records.SearchResult, error)

	// Delete removes record from vector store
	Delete(ctx context.Context, recID string) error

	// Close closes the vector store connection
	Close() error
}
