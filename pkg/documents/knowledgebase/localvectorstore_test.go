package knowledgebase

import (
	"context"
	"testing"

	"github.com/kazemisoroush/assistant/pkg/documents"
)

func TestLocalVectorStore_Index(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
	doc := &documents.Document{
		ID:          "doc1",
		Title:       "Go Programming",
		Content:     "Go is a great programming language",
		Description: "Introduction to Go",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, doc)

	// Assert
	if err != nil {
		t.Errorf("Index() error = %v, want nil", err)
	}
}

func TestLocalVectorStore_Index_MissingID(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
	doc := &documents.Document{
		Title:   "Go Programming",
		Content: "Go is a great programming language",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, doc)

	// Assert
	if err == nil {
		t.Error("Index() error = nil, want error for missing ID")
	}
}

func TestLocalVectorStore_Search(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
	doc := &documents.Document{
		ID:      "doc1",
		Title:   "Go Programming",
		Content: "Go is a great programming language for building scalable applications",
	}
	ctx := context.Background()
	if err := store.Index(ctx, doc); err != nil {
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
	if len(results) > 0 && results[0].Document.ID != "doc1" {
		t.Errorf("Search() returned document ID = %s, want doc1", results[0].Document.ID)
	}
}

func TestLocalVectorStore_Search_EmptyStore(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
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

func TestLocalVectorStore_Delete(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
	doc := &documents.Document{
		ID:      "doc1",
		Title:   "Test Document",
		Content: "Test content",
	}
	ctx := context.Background()
	if err := store.Index(ctx, doc); err != nil {
		t.Fatalf("Index() failed: %v", err)
	}

	// Act
	err := store.Delete(ctx, "doc1")

	// Assert
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}

	// Verify document is deleted
	results, err := store.Search(ctx, "test", 10)
	if err != nil {
		t.Errorf("Search() after Delete() error = %v, want nil", err)
	}
	if len(results) != 0 {
		t.Errorf("After Delete(), Search() returned %d results, want 0", len(results))
	}
}

func TestLocalVectorStore_Delete_NotFound(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()
	ctx := context.Background()

	// Act
	err := store.Delete(ctx, "nonexistent")

	// Assert
	if err == nil {
		t.Error("Delete() error = nil, want error for nonexistent document")
	}
}

func TestLocalVectorStore_Close(t *testing.T) {
	// Arrange
	store := NewLocalVectorStore()

	// Act
	err := store.Close()

	// Assert
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}
