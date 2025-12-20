package handler

import (
	"context"

	service "github.com/kazemisoroush/assistant/pkg/records/service"
)

// QueryHandler handles searching for records.
type QueryHandler struct {
	service service.Service
}

// NewQueryHandler creates a new search handler.
func NewQueryHandler(service service.Service) Handler {
	return &QueryHandler{
		service: service,
	}
}

// Handle implements Handler.
func (l QueryHandler) Handle(_ context.Context, _ Request) (Response, error) {
	return Response{}, nil
}
