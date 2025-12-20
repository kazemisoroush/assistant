// Package records implements the record service.
package records

import (
	"context"
	"fmt"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

// Service defines operations for record management
//
//go:generate mockgen -destination=./mocks/mock_service.go -mock_names=Service=MockService -package=mocks . Service
type Service interface {
	// Ingest processes and stores a record
	Ingest(ctx context.Context, rec records.Record) error

	// Search performs semantic search with optional metadata filters
	// For now this is basic keyword search, will be enhanced with vector search
	Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]records.SearchResult, error)

	// GetByID retrieves a record by its ID
	GetByID(ctx context.Context, id string) (records.Record, error)

	// List returns all records with optional type filter
	List(ctx context.Context, recType records.RecordType) ([]records.Record, error)

	// Update updates an existing record
	Update(ctx context.Context, rec records.Record) error

	// Delete removes a record
	Delete(ctx context.Context, id string) error
}

// RecordService implements the Service interface
type RecordService struct {
	storage       storage.Storage
	vectorStorage knowledgebase.VectorStorage
}

// NewRecordService creates a new record service
// vectorStorage can be nil if semantic search is not needed
func NewRecordService(storage storage.Storage, vectorStorage knowledgebase.VectorStorage) Service {
	return &RecordService{
		storage:       storage,
		vectorStorage: vectorStorage,
	}
}

// Ingest processes and stores a record
func (s *RecordService) Ingest(ctx context.Context, rec records.Record) error {
	// Store the record
	if err := s.storage.Store(ctx, rec); err != nil {
		return fmt.Errorf("failed to store record: %w", err)
	}

	// Index in vector store for semantic search
	if err := s.vectorStorage.Index(ctx, rec); err != nil {
		return fmt.Errorf("failed to index record: %w", err)
	}

	return nil
}

// Update updates an existing record
func (s *RecordService) Update(ctx context.Context, rec records.Record) error {
	if err := s.storage.Update(ctx, rec); err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	// Update in vector store (reindex with new content)
	if err := s.vectorStorage.Index(ctx, rec); err != nil {
		return fmt.Errorf("failed to reindex record: %w", err)
	}

	return nil
}

// Delete removes a record
func (s *RecordService) Delete(ctx context.Context, id string) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	// Delete from vector store
	if err := s.vectorStorage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete from vector store: %w", err)
	}

	return nil
}

// Search performs search with optional metadata filters
func (s *RecordService) Search(_ context.Context, _ string, _ map[string]interface{}, _ int) ([]records.SearchResult, error) {
	panic("not implemented")
}

// GetByID retrieves a record by its ID
func (s *RecordService) GetByID(ctx context.Context, id string) (records.Record, error) {
	return s.storage.Get(ctx, id)
}

// List returns all records with optional type filter
func (s *RecordService) List(ctx context.Context, recType records.RecordType) ([]records.Record, error) {
	return s.storage.List(ctx, recType)
}
