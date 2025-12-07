package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kazemisoroush/assistant/pkg/documents"
)

// LocalStorage implements document storage using local filesystem + JSON
// Documents are stored as individual JSON files with in-memory caching
type LocalStorage struct {
	basePath string
	mu       sync.RWMutex
	docs     map[string]*documents.Document // In-memory cache
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (Storage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	ls := &LocalStorage{
		basePath: basePath,
		docs:     make(map[string]*documents.Document),
	}

	// Load existing documents
	if err := ls.loadDocuments(); err != nil {
		return nil, fmt.Errorf("failed to load existing documents: %w", err)
	}

	return ls, nil
}

// Store saves a document to disk and cache
func (ls *LocalStorage) Store(_ context.Context, doc *documents.Document) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Update timestamp
	doc.UpdatedAt = time.Now()

	// Save to disk
	docPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", doc.ID))
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	if err := os.WriteFile(docPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	// Update cache
	ls.docs[doc.ID] = doc
	return nil
}

// Get retrieves a document by ID
func (ls *LocalStorage) Get(_ context.Context, id string) (*documents.Document, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	doc, exists := ls.docs[id]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	return doc, nil
}

// List returns all documents with optional type filter
func (ls *LocalStorage) List(_ context.Context, docType documents.DocumentType) ([]*documents.Document, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	var result []*documents.Document
	for _, doc := range ls.docs {
		if docType == "" || doc.Type == docType {
			result = append(result, doc)
		}
	}

	return result, nil
}

// Update updates an existing document
func (ls *LocalStorage) Update(_ context.Context, doc *documents.Document) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Check if document exists
	if _, exists := ls.docs[doc.ID]; !exists {
		return fmt.Errorf("document not found: %s", doc.ID)
	}

	// Update timestamp
	doc.UpdatedAt = time.Now()

	// Save to disk
	docPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", doc.ID))
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	if err := os.WriteFile(docPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}

	// Update cache
	ls.docs[doc.ID] = doc
	return nil
}

// Delete removes a document
func (ls *LocalStorage) Delete(_ context.Context, id string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Check if document exists
	if _, exists := ls.docs[id]; !exists {
		return fmt.Errorf("document not found: %s", id)
	}

	// Delete from disk
	docPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", id))
	if err := os.Remove(docPath); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Remove from cache
	delete(ls.docs, id)
	return nil
}

// Search performs basic keyword search across documents
// This is a simple implementation that will be enhanced with vector search later
func (ls *LocalStorage) Search(_ context.Context, query string, filters map[string]interface{}, limit int) ([]documents.SearchResult, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	var results []documents.SearchResult
	queryLower := strings.ToLower(query)

	for _, doc := range ls.docs {
		// Apply filters first
		if !matchesFilters(doc, filters) {
			continue
		}

		score := calculateSearchScore(doc, queryLower)
		if score > 0 {
			results = append(results, documents.SearchResult{
				Document: *doc,
				Score:    score,
			})
		}
	}

	// Sort by score (simple bubble sort for now)
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// Apply limit
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// calculateSearchScore computes a relevance score based on keyword matching
func calculateSearchScore(doc *documents.Document, queryLower string) float64 {
	score := 0.0
	contentLower := strings.ToLower(doc.Content)
	titleLower := strings.ToLower(doc.Title)
	descLower := strings.ToLower(doc.Description)

	if strings.Contains(titleLower, queryLower) {
		score += 0.5
	}
	if strings.Contains(descLower, queryLower) {
		score += 0.3
	}
	if strings.Contains(contentLower, queryLower) {
		score += 0.2
	}

	return score
}

// loadDocuments loads all documents from disk into memory
func (ls *LocalStorage) loadDocuments() error {
	entries, err := os.ReadDir(ls.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, that's okay
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		docPath := filepath.Join(ls.basePath, entry.Name())
		data, err := os.ReadFile(docPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		var doc documents.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", entry.Name(), err)
		}

		ls.docs[doc.ID] = &doc
	}

	return nil
}

// matchesFilters checks if a document matches the given filters
func matchesFilters(doc *documents.Document, filters map[string]interface{}) bool {
	if len(filters) == 0 {
		return true
	}

	for key, value := range filters {
		switch key {
		case "type":
			if !matchesTypeFilter(doc, value) {
				return false
			}
		case "tag":
			if !matchesTagFilter(doc, value) {
				return false
			}
		}
	}

	return true
}

func matchesTypeFilter(doc *documents.Document, value interface{}) bool {
	strVal, ok := value.(string)
	return !ok || doc.Type == documents.DocumentType(strVal)
}

func matchesTagFilter(doc *documents.Document, value interface{}) bool {
	tagValue, ok := value.(string)
	if !ok {
		return true
	}
	for _, tag := range doc.Tags {
		if tag == tagValue {
			return true
		}
	}
	return false
}
