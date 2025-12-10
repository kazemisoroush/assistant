// Package records provides record management functionality.
package records

import (
	"context"
)

// Service defines operations for record management
//
//go:generate mockgen -destination=./mocks/mock_service.go -mock_names=Service=MockService -package=mocks . Service
type Service interface {
	// Ingest processes and stores a record
	Ingest(ctx context.Context, rec *Record) error

	// Search performs semantic search with optional metadata filters
	// For now this is basic keyword search, will be enhanced with vector search
	Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]SearchResult, error)

	// GetByID retrieves a record by its ID
	GetByID(ctx context.Context, id string) (*Record, error)

	// List returns all records with optional type filter
	List(ctx context.Context, recType RecordType) ([]*Record, error)

	// Update updates an existing record
	Update(ctx context.Context, rec *Record) error

	// Delete removes a record
	Delete(ctx context.Context, id string) error
}
