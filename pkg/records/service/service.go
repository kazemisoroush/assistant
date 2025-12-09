// Package records implements the record service.
package records

import (
	"context"
	"fmt"
	"strings"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

// RecordService implements the Service interface
type RecordService struct {
	storage     storage.Storage
	vectorStore knowledgebase.VectorStore // Vector store for semantic search
}

// NewRecordService creates a new record service
// vectorStore can be nil if semantic search is not needed
func NewRecordService(storage storage.Storage, vectorStore knowledgebase.VectorStore) records.Service {
	return &RecordService{
		storage:     storage,
		vectorStore: vectorStore,
	}
}

// Ingest processes and stores a record
func (s *RecordService) Ingest(ctx context.Context, rec *records.Record) error {
	// Validate record
	if rec.ID == "" {
		return fmt.Errorf("record ID is required")
	}
	if rec.Type == "" {
		return fmt.Errorf("record type is required")
	}
	if rec.Content == "" && rec.FilePath == "" {
		return fmt.Errorf("record must have either content or file path")
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
	if s.vectorStore != nil {
		if err := s.vectorStore.Index(ctx, rec); err != nil {
			return fmt.Errorf("failed to index record: %w", err)
		}
	}

	return nil
}

// Search performs search with optional metadata filters
func (s *RecordService) Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]records.SearchResult, error) {
	// Use vector store for semantic search if available
	if s.vectorStore != nil {
		results, err := s.vectorStore.Search(ctx, query, limit)
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
func (s *RecordService) GetByID(ctx context.Context, id string) (*records.Record, error) {
	return s.storage.Get(ctx, id)
}

// List returns all records with optional type filter
func (s *RecordService) List(ctx context.Context, recType records.RecordType) ([]*records.Record, error) {
	return s.storage.List(ctx, recType)
}

// Update updates an existing record
func (s *RecordService) Update(ctx context.Context, rec *records.Record) error {
	if rec.ID == "" {
		return fmt.Errorf("record ID is required")
	}

	if err := s.storage.Update(ctx, rec); err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	// Update in vector store (reindex with new content)
	if s.vectorStore != nil {
		if err := s.vectorStore.Index(ctx, rec); err != nil {
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
	if s.vectorStore != nil {
		if err := s.vectorStore.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete from vector store: %w", err)
		}
	}

	return nil
}

// ExtractTextFromFile is a helper function to extract text content from various file types
// For now, it just reads plain text. Later we can add PDF, DOCX, image OCR support
func ExtractTextFromFile(_ string) (string, error) {
	// TODO: Implement based on file type
	// - .txt: read directly
	// - .pdf: use pdf library
	// - .docx: use docx library
	// - .jpg, .png: use OCR
	return "", fmt.Errorf("not implemented yet")
}

// NormalizeContent performs basic text normalization
func NormalizeContent(content string) string {
	// Trim whitespace
	content = strings.TrimSpace(content)
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	return content
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
