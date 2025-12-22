// Package source provides interfaces and implementations for data sources.
package source

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/kazemisoroush/assistant/pkg/records/extractor"
)

// LocalSource reads files from a local directory structure
type LocalSource struct {
	extractor extractor.ContentExtractor
	basePath  string
}

// NewLocalSource creates a new local file source
func NewLocalSource(extractor extractor.ContentExtractor, basePath string) Source {
	return &LocalSource{
		extractor: extractor,
		basePath:  basePath,
	}
}

// Name returns the source name
func (ls *LocalSource) Name() string {
	return "local"
}

// Scrape reads files from the local directory structure
func (ls *LocalSource) Scrape(ctx context.Context) (<-chan records.Record, <-chan error) {
	recordChan := make(chan records.Record)
	errChan := make(chan error, 1)

	go func() {
		defer close(recordChan)
		defer close(errChan)

		err := filepath.WalkDir(ls.basePath, func(path string, d fs.DirEntry, err error) error {
			// Check for context cancellation
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if err != nil {
				return err
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				errChan <- fmt.Errorf("failed to read file %s: %w", path, err)
				return nil // Continue processing other files
			}

			record, err := ls.extractor.Extract(ctx, string(content))
			if err != nil {
				errChan <- fmt.Errorf("failed to extract record from file %s: %w", path, err)
				return nil // Continue processing other files
			}

			recordChan <- record
			return nil
		})

		if err != nil {
			errChan <- fmt.Errorf("failed to walk directory: %w", err)
		}
	}()

	return recordChan, errChan
}
