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

	// TODO: Here we want to implement a logic to walk through the base path and scrape all files.
	// For each file found we need to pass it to a new interface named "ItemExtractor" or a better name if you can think of.
	// A particular implementation of ItemExtractor probably file extractor or a better name will be responsible to read and OCR for non-text files.
	// For OCR extraction, record type is determined by the file OCR result. I.e. there should be probably a separate method for type detection.
	// After extracting the content, we need to create a records.Record object and send it to recordChan.
	// In case of any error during the process, we should send the error to errChan.

	return recordChan, errChan
}
