package storage

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/documents"
)

// Storage defines the persistence layer interface
type Storage interface {
	// Store saves a document
	Store(ctx context.Context, doc *documents.Document) error

	// Get retrieves a document by ID
	Get(ctx context.Context, id string) (*documents.Document, error)

	// List returns all documents with optional type filter
	List(ctx context.Context, docType documents.DocumentType) ([]*documents.Document, error)

	// Update updates an existing document
	Update(ctx context.Context, doc *documents.Document) error

	// Delete removes a document
	Delete(ctx context.Context, id string) error
}
