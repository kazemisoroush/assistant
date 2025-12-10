// Package source provides interfaces and implementations for data sources.
package source

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kazemisoroush/assistant/pkg/records"
)

// LocalSource reads files from a local directory structure
type LocalSource struct {
	name     string
	basePath string
	enabled  bool
}

// NewLocalSource creates a new local file source
func NewLocalSource(name, basePath string, enabled bool) Source {
	return &LocalSource{
		name:     name,
		basePath: basePath,
		enabled:  enabled,
	}
}

// Name returns the source name
func (ls *LocalSource) Name() string {
	return ls.name
}

// IsEnabled returns whether this source is enabled
func (ls *LocalSource) IsEnabled() bool {
	return ls.enabled
}

// Scrape reads files from the local directory structure
// TODO: WOW What the hell? This is wrong in so many levels. Why is this returning channels?
func (ls *LocalSource) Scrape(ctx context.Context) (<-chan *records.Record, <-chan error) {
	recordChan := make(chan *records.Record)
	errChan := make(chan error, 1)

	go func() {
		defer close(recordChan)
		defer close(errChan)

		// Check if directory exists
		if _, err := os.Stat(ls.basePath); os.IsNotExist(err) {
			errChan <- fmt.Errorf("base path does not exist: %s", ls.basePath)
			return
		}

		// Walk through the directory structure
		err := filepath.Walk(ls.basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Skip non-text files
			if !strings.HasSuffix(path, ".txt") {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				errChan <- fmt.Errorf("failed to read file %s: %w", path, err)
				return nil // Continue with other files
			}

			// Determine record type from directory structure
			relPath, _ := filepath.Rel(ls.basePath, path)
			recordType := determineRecordType(relPath)

			// Create record
			rec := &records.Record{
				ID:        uuid.New().String(),
				Type:      recordType,
				FilePath:  path,
				FileName:  info.Name(),
				Title:     generateTitle(info.Name()),
				Content:   string(content),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Metadata:  extractMetadata(recordType, string(content)),
				Tags:      []string{"scraped", "local"},
			}

			select {
			case recordChan <- rec:
			case <-ctx.Done():
				return ctx.Err()
			}

			return nil
		})

		if err != nil && err != context.Canceled {
			errChan <- fmt.Errorf("error walking directory: %w", err)
		}
	}()

	return recordChan, errChan
}

// determineRecordType maps directory names to record types
// TODO: The record type is determined contextual. You can't expect to know the record type based on the relative path. Instead you should ask LLM and let it decide what is the record type...
func determineRecordType(relPath string) records.RecordType {
	parts := strings.Split(relPath, string(os.PathSeparator))
	if len(parts) == 0 {
		return records.RecordTypeOther
	}

	// Use the first directory name to determine type
	dirName := parts[0]

	switch strings.ToLower(dirName) {
	case "health_visit", "health-visit":
		return records.RecordTypeHealthVisit
	case "health_test", "health-test":
		return records.RecordTypeHealthTest
	case "health_lab", "health-lab":
		return records.RecordTypeHealthLab
	case "receipts", "receipt":
		return records.RecordTypeReceipt
	case "insurance":
		return records.RecordTypeInsurance
	case "id", "identification":
		return records.RecordTypeID
	case "travel":
		return records.RecordTypeTravel
	case "work_contract", "work-contract", "work":
		return records.RecordTypeWorkContract
	case "tax", "taxes":
		return records.RecordTypeTax
	case "car", "vehicle":
		return records.RecordTypeCar
	case "home", "house", "property":
		return records.RecordTypeHome
	case "visa", "visas":
		return records.RecordTypeVisa
	default:
		return records.RecordTypeOther
	}
}

// generateTitle creates a human-readable title from filename
func generateTitle(filename string) string {
	// Remove extension
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Replace underscores and hyphens with spaces
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	// Capitalize first letter of each word
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// extractMetadata extracts basic metadata from content
func extractMetadata(recordType records.RecordType, content string) map[string]interface{} {
	metadata := make(map[string]interface{})

	// Add record type
	metadata["record_type"] = string(recordType)

	// Add content length
	metadata["content_length"] = len(content)

	// Add word count
	words := strings.Fields(content)
	metadata["word_count"] = len(words)

	// Type-specific metadata extraction (basic)
	switch recordType {
	case records.RecordTypeReceipt:
		// Look for common receipt patterns
		if strings.Contains(strings.ToLower(content), "total") {
			metadata["has_total"] = true
		}
		if strings.Contains(strings.ToLower(content), "date") {
			metadata["has_date"] = true
		}
	case records.RecordTypeHealthVisit:
		// Look for medical visit patterns
		if strings.Contains(strings.ToLower(content), "doctor") {
			metadata["has_doctor_info"] = true
		}
		if strings.Contains(strings.ToLower(content), "diagnosis") {
			metadata["has_diagnosis"] = true
		}
	}

	return metadata
}
