package ingestor

import (
	"context"
	"fmt"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

// RecordIngestor is an implementation of the Ingestor interface.
type RecordIngestor struct {
	storage       storage.Storage
	vectorStorage knowledgebase.VectorStorage
}

// NewRecordIngestor creates a new instance of RecordIngestor.
func NewRecordIngestor(storage storage.Storage, vectorStorage knowledgebase.VectorStorage) Ingestor {
	return &RecordIngestor{
		storage:       storage,
		vectorStorage: vectorStorage,
	}
}

// Ingest processes and stores a record without checking for existing records
// Ingest processes and stores a record (upsert behavior)
func (s *RecordIngestor) Ingest(ctx context.Context, record records.Record) error {
	// Check if record exists
	_, err := s.storage.Get(ctx, record.ID)
	if err == nil {
		// Record exists, delete from both storage and vector store
		if err := s.storage.Delete(ctx, record.ID); err != nil {
			return fmt.Errorf("failed to delete existing record from storage: %w", err)
		}
		if err := s.vectorStorage.Delete(ctx, record.ID); err != nil {
			return fmt.Errorf("failed to delete existing record from vector store: %w", err)
		}
	}

	// Store the record
	if err := s.storage.Store(ctx, record); err != nil {
		return fmt.Errorf("failed to store record: %w", err)
	}

	// Index in vector store for semantic search
	if err := s.vectorStorage.Index(ctx, record); err != nil {
		return fmt.Errorf("failed to index record: %w", err)
	}

	return nil
}

// Delete removes a record
func (s *RecordIngestor) Delete(ctx context.Context, id string) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	// Delete from vector store
	if err := s.vectorStorage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete from vector store: %w", err)
	}

	return nil
}
