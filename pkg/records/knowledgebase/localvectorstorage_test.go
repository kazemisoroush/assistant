package knowledgebase

import (
	"context"
	"testing"

	"github.com/kazemisoroush/assistant/pkg/records"
)

func TestLocalVectorStorage_Index(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := &records.Record{
		ID:      "rec1",
		Content: "Go is a great programming language",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, rec)

	// Assert
	if err != nil {
		t.Errorf("Index() error = %v, want nil", err)
	}
}

func TestLocalVectorStorage_Index_MissingID(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := &records.Record{
		Content: "Go is a great programming language",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, rec)

	// Assert
	if err == nil {
		t.Error("Index() error = nil, want error for missing ID")
	}
}

func TestLocalVectorStorage_Search(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := &records.Record{
		ID:      "rec1",
		Content: "Go is a great programming language for building scalable applications",
	}
	ctx := context.Background()
	if err := store.Index(ctx, rec); err != nil {
		t.Fatalf("Index() failed: %v", err)
	}

	// Act
	results, err := store.Search(ctx, "programming language", 10)

	// Assert
	if err != nil {
		t.Errorf("Search() error = %v, want nil", err)
	}
	if len(results) == 0 {
		t.Error("Search() returned no results, want at least 1")
	}
	if len(results) > 0 && results[0].Record.ID != "rec1" {
		t.Errorf("Search() returned record ID = %s, want rec1", results[0].Record.ID)
	}
}

func TestLocalVectorStorage_Search_EmptyStore(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	ctx := context.Background()

	// Act
	results, err := store.Search(ctx, "test query", 10)

	// Assert
	if err != nil {
		t.Errorf("Search() error = %v, want nil", err)
	}
	if len(results) != 0 {
		t.Errorf("Search() returned %d results, want 0", len(results))
	}
}

func TestLocalVectorStorage_Delete(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := &records.Record{
		ID:      "rec1",
		Content: "Test content",
	}
	ctx := context.Background()
	if err := store.Index(ctx, rec); err != nil {
		t.Fatalf("Index() failed: %v", err)
	}

	// Act
	err := store.Delete(ctx, "rec1")

	// Assert
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}

	// Verify record is deleted
	results, err := store.Search(ctx, "test", 10)
	if err != nil {
		t.Errorf("Search() after Delete() error = %v, want nil", err)
	}
	if len(results) != 0 {
		t.Errorf("After Delete(), Search() returned %d results, want 0", len(results))
	}
}

func TestLocalVectorStorage_Delete_NotFound(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	ctx := context.Background()

	// Act
	err := store.Delete(ctx, "nonexistent")

	// Assert
	if err == nil {
		t.Error("Delete() error = nil, want error for nonexistent record")
	}
}

func TestLocalVectorStorage_Close(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()

	// Act
	err := store.Close()

	// Assert
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}
