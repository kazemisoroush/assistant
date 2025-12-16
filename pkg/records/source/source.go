// Package source provides interfaces and implementations for data sources.
package source

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// Source represents a source of records that can be scraped/ingested
//
//go:generate mockgen -destination=./mocks/mock_source.go -mock_names=Source=MockSource -package=mocks . Source
type Source interface {
	// Name returns the name/identifier of this source
	Name() string

	// Scrape retrieves records from this source
	// Returns a channel of records and an error channel
	Scrape(ctx context.Context) (<-chan records.Record, <-chan error)
}
