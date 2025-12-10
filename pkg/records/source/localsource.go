// Package source provides interfaces and implementations for data sources.
package source

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// LocalSource reads files from a local directory structure
type LocalSource struct {
	basePath string
}

// NewLocalSource creates a new local file source
func NewLocalSource(basePath string) Source {
	return &LocalSource{
		basePath: basePath,
	}
}

// Name returns the source name
func (ls *LocalSource) Name() string {
	return "local"
}

// Scrape reads files from the local directory structure
func (ls *LocalSource) Scrape(_ context.Context) (<-chan *records.Record, <-chan error) {
	recordChan := make(chan *records.Record)
	errChan := make(chan error, 1)

	return recordChan, errChan
}
