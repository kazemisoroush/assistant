package documents

import (
	"context"
)

// Service defines operations for document management
type Service interface {
	// Ingest processes and stores a document
	Ingest(ctx context.Context, doc *Document) error

	// Search performs semantic search with optional metadata filters
	// For now this is basic keyword search, will be enhanced with vector search
	Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]SearchResult, error)

	// GetByID retrieves a document by its ID
	GetByID(ctx context.Context, id string) (*Document, error)

	// List returns all documents with optional type filter
	List(ctx context.Context, docType DocumentType) ([]*Document, error)

	// Update updates an existing document
	Update(ctx context.Context, doc *Document) error

	// Delete removes a document
	Delete(ctx context.Context, id string) error
}
