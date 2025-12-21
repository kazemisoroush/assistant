package knowledgebase

import (
	"context"
	"testing"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalVectorStorage_Index(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := records.Record{
		ID:      "rec1",
		Content: "Go is a great programming language",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, rec)

	// Assert
	require.NoError(t, err, "Index() error should be nil")
}

func TestLocalVectorStorage_Index_MissingID(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := records.Record{
		Content: "Go is a great programming language",
	}
	ctx := context.Background()

	// Act
	err := store.Index(ctx, rec)

	// Assert
	require.Error(t, err, "Index() error should not be nil for missing ID")
}

func TestLocalVectorStorage_Search(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := records.Record{
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
	require.NoError(t, err, "Search() error should be nil")
	assert.Greater(t, len(results), 0, "Search() should return at least one result")
	assert.Equal(t, "rec1", results[0].Record.ID, "Search() should return the indexed record")
}

func TestLocalVectorStorage_Search_EmptyStore(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	ctx := context.Background()

	// Act
	results, err := store.Search(ctx, "test query", 10)

	// Assert
	require.NoError(t, err, "Search() error should be nil")
	assert.Equal(t, 0, len(results), "Search() should return no results")
}

func TestLocalVectorStorage_Delete(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	rec := records.Record{
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
	require.NoError(t, err, "Delete() error should be nil")

	// Verify record is deleted
	results, err := store.Search(ctx, "test", 10)
	require.NoError(t, err, "Search() after Delete() error should be nil")
	assert.Equal(t, 0, len(results), "After Delete(), Search() should return no results")
}

func TestLocalVectorStorage_Delete_NotFound(t *testing.T) {
	// Arrange
	store := NewLocalVectorStorage()
	ctx := context.Background()

	// Act
	err := store.Delete(ctx, "nonexistent")

	// Assert
	require.Error(t, err, "Delete() error should not be nil for nonexistent record")
}
