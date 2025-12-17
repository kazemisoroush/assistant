package handler

import (
	"context"

	service "github.com/kazemisoroush/assistant/pkg/records/service"
	"github.com/kazemisoroush/assistant/pkg/records/source"
)

// LocalScraperHandler handles scraping records from local sources.
type LocalScraperHandler struct {
	service service.Service
	sources []source.Source
}

// NewLocalScraperHandler creates a new local scraper handler.
func NewLocalScraperHandler(service service.Service, sources []source.Source) Handler {
	return &LocalScraperHandler{
		service: service,
		sources: sources,
	}
}

// Handle implements Handler.
func (l LocalScraperHandler) Handle(_ context.Context, _ Request) (Response, error) {
	return Response{}, nil
}
