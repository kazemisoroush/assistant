package handler

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records/ingestor"
	"github.com/kazemisoroush/assistant/pkg/records/source"
)

// LocalScraperHandler handles scraping records from local sources.
type LocalScraperHandler struct {
	ingestor ingestor.Ingestor
	sources  []source.Source
}

// NewLocalScraperHandler creates a new local scraper handler.
func NewLocalScraperHandler(ingestor ingestor.Ingestor, sources []source.Source) Handler {
	return &LocalScraperHandler{
		ingestor: ingestor,
		sources:  sources,
	}
}

// Handle implements Handler.
func (l LocalScraperHandler) Handle(_ context.Context, _ Request) (Response, error) {
	return Response{}, nil
}
