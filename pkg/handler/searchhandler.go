package handler

import (
	"context"
	"fmt"
	"os"

	service "github.com/kazemisoroush/assistant/pkg/records/service"
)

// SearchHandler handles searching for records.
type SearchHandler struct {
	service service.Service
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(service service.Service) Handler {
	return &SearchHandler{
		service: service,
	}
}

// Handle implements Handler.
func (l SearchHandler) Handle(ctx context.Context) {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Error: search query is required\n\n")
		fmt.Println("Usage: assistant search <query>")
		fmt.Println("Example: assistant search \"health visit doctor\"")
		os.Exit(1)
	}

	query := os.Args[2]

	fmt.Printf("🔎 Searching for: %s\n", query)

	// Reindex existing records for search (vector store is in-memory)
	fmt.Println("📚 Loading and indexing existing records...")

	// Get all records from storage
	allRecords, err := l.service.List(ctx, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list records: %v\n", err)
		os.Exit(1)
	}

	// Reindex each record
	count := 0
	for _, rec := range allRecords {
		if err := l.service.Ingest(ctx, rec); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to reindex record %s: %v\n", rec.ID, err)
		} else {
			count++
		}
	}

	fmt.Printf("Indexed %d records\n", count)
	fmt.Println()

	// Perform search
	results, err := l.service.Search(ctx, query, nil, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search failed: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("Found %d results:\n\n", len(results))

	for _, result := range results {
		rec := result.Record
		fmt.Printf("   Type: %s\n", rec.Type)
		fmt.Printf("   ID: %s\n", rec.ID)
		if result.Score > 0 {
			fmt.Printf("   Relevance: %.2f\n", result.Score)
		}

		// Show snippet of content
		content := rec.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		fmt.Printf("   Preview: %s\n", content)
		fmt.Println()
	}
}
