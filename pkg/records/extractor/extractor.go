// Package extractor provides interfaces and implementations for extracting and classifying records from various content types.
package extractor

import "github.com/kazemisoroush/assistant/pkg/records"

// ContentExtractor defines an interface for extracting records from raw content.
//
//go:generate mockgen -destination=./mocks/mock_extractor.go -mock_names=ContentExtractor=MockContentExtractor -package=mocks . ContentExtractor
type ContentExtractor interface {
	// Extract processes raw content and returns a Record
	Extract(rawContent string) (records.Record, error)
}

// TypeExtractor defines an interface for classifying record types from text content.
type TypeExtractor interface {
	// GetType classifies the record type based on raw content
	GetType(textContent string) records.RecordType
}
