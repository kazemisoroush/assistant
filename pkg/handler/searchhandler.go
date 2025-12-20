package handler

import (
	"context"

	"github.com/kazemisoroush/assistant/pkg/records/ingestor"
)

// QueryHandler handles searching for records.
type QueryHandler struct {
	ingestor ingestor.Ingestor
}

// NewQueryHandler creates a new search handler.
func NewQueryHandler(ingestor ingestor.Ingestor) Handler {
	return &QueryHandler{
		ingestor: ingestor,
	}
}

// Handle implements Handler.
func (l QueryHandler) Handle(_ context.Context, _ Request) (Response, error) {
	return Response{}, nil
}
