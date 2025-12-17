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
	// Validate record
	if rec.ID == "" {
		return fmt.Errorf("record ID is required")
	}
	if rec.Type == "" {
		return fmt.Errorf("record type is required")
	}
	if rec.Content == "" {
		return fmt.Errorf("record must have content")
	}

	// Initialize metadata map if nil
	if rec.Metadata == nil {
		rec.Metadata = make(map[string]interface{})
	}

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

// Search performs search with optional metadata filters
func (s *RecordService) Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]records.SearchResult, error) {
	// Use vector store for semantic search if available
	if s.vectorStorage != nil {
		results, err := s.vectorStorage.Search(ctx, query, limit)
		if err != nil {
			return nil, fmt.Errorf("vector search failed: %w", err)
		}

		// Apply metadata filters if provided
		if len(filters) > 0 {
			results = applyFilters(results, filters)
		}

		return results, nil
	}

	// Fallback to basic keyword search from storage
	if localStorage, ok := s.storage.(interface {
		Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]records.SearchResult, error)
	}); ok {
		return localStorage.Search(ctx, query, filters, limit)
	}

	return nil, fmt.Errorf("search not supported by current storage implementation")
}

// GetByID retrieves a record by its ID
func (s *RecordService) GetByID(ctx context.Context, id string) (records.Record, error) {
	return s.storage.Get(ctx, id)
}

// List returns all records with optional type filter
func (s *RecordService) List(ctx context.Context, recType records.RecordType) ([]records.Record, error) {
	return s.storage.List(ctx, recType)
}

// Update updates an existing record
func (s *RecordService) Update(ctx context.Context, rec records.Record) error {
	if rec.ID == "" {
		return fmt.Errorf("record ID is required")
	}

	if err := s.storage.Update(ctx, rec); err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	// Update in vector store (reindex with new content)
	if s.vectorStorage != nil {
		if err := s.vectorStorage.Index(ctx, rec); err != nil {
			return fmt.Errorf("failed to reindex record: %w", err)
		}
	}

	return nil
}

// Delete removes a record
func (s *RecordService) Delete(ctx context.Context, id string) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	// Delete from vector store
	if s.vectorStorage != nil {
		if err := s.vectorStorage.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete from vector store: %w", err)
		}
	}

	return nil
}

// applyFilters filters search results based on metadata criteria
func applyFilters(results []records.SearchResult, filters map[string]interface{}) []records.SearchResult {
	if len(filters) == 0 {
		return results
	}

	filtered := make([]records.SearchResult, 0, len(results))
	for _, result := range results {
		if matchesFilters(&result.Record, filters) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

// matchesFilters checks if a record matches all filter criteria
func matchesFilters(rec *records.Record, filters map[string]interface{}) bool {
	for key, value := range filters {
		switch key {
		case "type":
			if rec.Type != records.RecordType(fmt.Sprint(value)) {
				return false
			}
		default:
			// Check in metadata
			if rec.Metadata == nil {
				return false
			}
			metaValue, exists := rec.Metadata[key]
			if !exists || metaValue != value {
				return false
			}
		}
	}
	return true
}
