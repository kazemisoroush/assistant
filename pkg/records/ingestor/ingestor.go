// Package ingestor implements the record ingestor service.
package ingestor

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// Ingestor defines operations for record management
//
//go:generate mockgen -destination=./mocks/mock_service.go -mock_names=Ingestor=MockService -package=mocks . Ingestor
type Ingestor interface {
	// Ingest processes and stores a record
	Ingest(ctx context.Context, record records.Record) error

	// Delete removes a record
	Delete(ctx context.Context, id string) error
}
