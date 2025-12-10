package storage

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// Storage defines the persistence layer interface
//
//go:generate mockgen -destination=./mocks/mock_storage.go -mock_names=Storage=MockStorage -package=mocks . Storage
type Storage interface {
	// Store saves a record
	Store(ctx context.Context, rec *records.Record) error

	// Get retrieves a record by ID
	Get(ctx context.Context, id string) (*records.Record, error)

	// List returns all records with optional type filter
	List(ctx context.Context, recType records.RecordType) ([]*records.Record, error)

	// Update updates an existing record
	Update(ctx context.Context, rec *records.Record) error

	// Delete removes a record
	Delete(ctx context.Context, id string) error
}
