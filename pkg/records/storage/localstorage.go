// Package storage implements storage backends for records.
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

	"github.com/kazemisoroush/assistant/pkg/records"
)

// LocalStorage implements record storage using local filesystem + JSON
// Records are stored as individual JSON files with in-memory caching
type LocalStorage struct {
	basePath string
	mu       sync.RWMutex
	recs     map[string]*records.Record // In-memory cache
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (Storage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	ls := &LocalStorage{
		basePath: basePath,
		recs:     make(map[string]*records.Record),
	}

	// Load existing records
	if err := ls.loadRecords(); err != nil {
		return nil, fmt.Errorf("failed to load existing records: %w", err)
	}

	return ls, nil
}

// Store saves a record to disk and cache
func (ls *LocalStorage) Store(_ context.Context, rec *records.Record) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Update timestamp
	rec.UpdatedAt = time.Now()

	// Save to disk
	recPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", rec.ID))
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	if err := os.WriteFile(recPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	// Update cache
	ls.recs[rec.ID] = rec
	return nil
}

// Get retrieves a record by ID
func (ls *LocalStorage) Get(_ context.Context, id string) (*records.Record, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	rec, exists := ls.recs[id]
	if !exists {
		return nil, fmt.Errorf("record not found: %s", id)
	}

	return rec, nil
}

// List returns all records with optional type filter
func (ls *LocalStorage) List(_ context.Context, recType records.RecordType) ([]*records.Record, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	var result []*records.Record
	for _, rec := range ls.recs {
		if recType == "" || rec.Type == recType {
			result = append(result, rec)
		}
	}

	return result, nil
}

// Update updates an existing record
func (ls *LocalStorage) Update(_ context.Context, rec *records.Record) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Check if record exists
	if _, exists := ls.recs[rec.ID]; !exists {
		return fmt.Errorf("record not found: %s", rec.ID)
	}

	// Update timestamp
	rec.UpdatedAt = time.Now()

	// Save to disk
	recPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", rec.ID))
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	if err := os.WriteFile(recPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	// Update cache
	ls.recs[rec.ID] = rec
	return nil
}

// Delete removes a record
func (ls *LocalStorage) Delete(_ context.Context, id string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Check if record exists
	if _, exists := ls.recs[id]; !exists {
		return fmt.Errorf("record not found: %s", id)
	}

	// Delete from disk
	recPath := filepath.Join(ls.basePath, fmt.Sprintf("%s.json", id))
	if err := os.Remove(recPath); err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	// Remove from cache
	delete(ls.recs, id)
	return nil
}

// Search performs basic keyword search across records
// This is a simple implementation that will be enhanced with vector search later
func (ls *LocalStorage) Search(_ context.Context, query string, filters map[string]interface{}, limit int) ([]records.SearchResult, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	var results []records.SearchResult
	queryLower := strings.ToLower(query)

	for _, rec := range ls.recs {
		// Apply filters first
		if !matchesFilters(rec, filters) {
			continue
		}

		score := calculateSearchScore(rec, queryLower)
		if score > 0 {
			results = append(results, records.SearchResult{
				Record: *rec,
				Score:  score,
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
func calculateSearchScore(rec *records.Record, queryLower string) float64 {
	score := 0.0
	contentLower := strings.ToLower(rec.Content)
	titleLower := strings.ToLower(rec.Title)
	descLower := strings.ToLower(rec.Description)

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

// loadRecords loads all records from disk into memory
func (ls *LocalStorage) loadRecords() error {
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

		recPath := filepath.Join(ls.basePath, entry.Name())
		data, err := os.ReadFile(recPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		var rec records.Record
		if err := json.Unmarshal(data, &rec); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", entry.Name(), err)
		}

		ls.recs[rec.ID] = &rec
	}

	return nil
}

// matchesFilters checks if a record matches the given filters
func matchesFilters(rec *records.Record, filters map[string]interface{}) bool {
	if len(filters) == 0 {
		return true
	}

	for key, value := range filters {
		switch key {
		case "type":
			if !matchesTypeFilter(rec, value) {
				return false
			}
		case "tag":
			if !matchesTagFilter(rec, value) {
				return false
			}
		}
	}

	return true
}

func matchesTypeFilter(rec *records.Record, value interface{}) bool {
	strVal, ok := value.(string)
	return !ok || rec.Type == records.RecordType(strVal)
}

func matchesTagFilter(rec *records.Record, value interface{}) bool {
	tagValue, ok := value.(string)
	if !ok {
		return true
	}
	for _, tag := range rec.Tags {
		if tag == tagValue {
			return true
		}
	}
	return false
}
