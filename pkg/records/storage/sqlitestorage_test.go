package storage

import (
	"context"
	"testing"
	"time"

	"github.com/kazemisoroush/assistant/pkg/records"
)

func setupTestDB(t *testing.T) (*SQLiteStorage, func()) {
	t.Helper()

	// Use in-memory database for testing
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("failed to create test storage: %v", err)
	}

	cleanup := func() {
		_ = storage.Close()
	}

	return storage, cleanup
}

func createTestRecord(id string, recType records.RecordType) records.Record {
	return records.Record{
		ID:        id,
		Type:      recType,
		Content:   "test content for " + id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"test_key": "test_value",
			"number":   42,
		},
		Tags: []string{"test", "tag1"},
	}
}

func TestStore(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rec := createTestRecord("test-id-1", records.RecordTypeReceipt)

	err := storage.Store(ctx, rec)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Verify it was stored by retrieving it
	retrieved, err := storage.Get(ctx, rec.ID)
	if err != nil {
		t.Fatalf("Get failed after Store: %v", err)
	}

	if retrieved.ID != rec.ID {
		t.Errorf("expected ID %s, got %s", rec.ID, retrieved.ID)
	}
	if retrieved.Type != rec.Type {
		t.Errorf("expected Type %s, got %s", rec.Type, retrieved.Type)
	}
	if retrieved.Content != rec.Content {
		t.Errorf("expected Content %s, got %s", rec.Content, retrieved.Content)
	}
}

func TestGet(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rec := createTestRecord("test-id-2", records.RecordTypeHealthVisit)

	// Store first
	if err := storage.Store(ctx, rec); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Test Get
	retrieved, err := storage.Get(ctx, rec.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.ID != rec.ID {
		t.Errorf("expected ID %s, got %s", rec.ID, retrieved.ID)
	}
	if retrieved.Metadata["test_key"] != "test_value" {
		t.Errorf("metadata not properly retrieved")
	}
}

func TestGet_NotFound(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	_, err := storage.Get(ctx, "non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent record, got nil")
	}
}

func TestList(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Store multiple records
	rec1 := createTestRecord("id-1", records.RecordTypeReceipt)
	rec2 := createTestRecord("id-2", records.RecordTypeReceipt)
	rec3 := createTestRecord("id-3", records.RecordTypeHealthVisit)

	for _, rec := range []records.Record{rec1, rec2, rec3} {
		if err := storage.Store(ctx, rec); err != nil {
			t.Fatalf("Store failed: %v", err)
		}
	}

	// List all records
	allRecords, err := storage.List(ctx, "")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(allRecords) != 3 {
		t.Errorf("expected 3 records, got %d", len(allRecords))
	}
}

func TestList_WithFilter(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Store multiple records
	rec1 := createTestRecord("id-1", records.RecordTypeReceipt)
	rec2 := createTestRecord("id-2", records.RecordTypeReceipt)
	rec3 := createTestRecord("id-3", records.RecordTypeHealthVisit)

	for _, rec := range []records.Record{rec1, rec2, rec3} {
		if err := storage.Store(ctx, rec); err != nil {
			t.Fatalf("Store failed: %v", err)
		}
	}

	// List only receipts
	receipts, err := storage.List(ctx, records.RecordTypeReceipt)
	if err != nil {
		t.Fatalf("List with filter failed: %v", err)
	}

	if len(receipts) != 2 {
		t.Errorf("expected 2 receipts, got %d", len(receipts))
	}

	for _, rec := range receipts {
		if rec.Type != records.RecordTypeReceipt {
			t.Errorf("expected type %s, got %s", records.RecordTypeReceipt, rec.Type)
		}
	}
}

func TestUpdate(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rec := createTestRecord("test-id-4", records.RecordTypeReceipt)

	// Store first
	if err := storage.Store(ctx, rec); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Update
	rec.Content = "updated content"
	rec.Type = records.RecordTypeHealthLab
	rec.UpdatedAt = time.Now()

	err := storage.Update(ctx, rec)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrieved, err := storage.Get(ctx, rec.ID)
	if err != nil {
		t.Fatalf("Get failed after Update: %v", err)
	}

	if retrieved.Content != "updated content" {
		t.Errorf("expected updated content, got %s", retrieved.Content)
	}
	if retrieved.Type != records.RecordTypeHealthLab {
		t.Errorf("expected type %s, got %s", records.RecordTypeHealthLab, retrieved.Type)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rec := createTestRecord("non-existent", records.RecordTypeReceipt)

	err := storage.Update(ctx, rec)
	if err == nil {
		t.Error("expected error for updating non-existent record, got nil")
	}
}

func TestDelete(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	rec := createTestRecord("test-id-5", records.RecordTypeReceipt)

	// Store first
	if err := storage.Store(ctx, rec); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Delete
	err := storage.Delete(ctx, rec.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = storage.Get(ctx, rec.ID)
	if err == nil {
		t.Error("expected error when getting deleted record, got nil")
	}
}

func TestDelete_NotFound(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	err := storage.Delete(ctx, "non-existent-id")
	if err == nil {
		t.Error("expected error for deleting non-existent record, got nil")
	}
}

func TestClose(t *testing.T) {
	storage, _ := setupTestDB(t)

	err := storage.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Attempting operations after close should fail
	ctx := context.Background()
	rec := createTestRecord("test-id-6", records.RecordTypeReceipt)

	err = storage.Store(ctx, rec)
	if err == nil {
		t.Error("expected error when using closed storage, got nil")
	}
}
