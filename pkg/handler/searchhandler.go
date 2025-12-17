package handler

import (
	"context"

	service "github.com/kazemisoroush/assistant/pkg/records/service"
)

// SearchHandler handles searching for records.
type SearchHandler struct {
	service service.Service
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(service service.Service) Handler {
	return &SearchHandler{
		service: service,
	}
}

// Handle implements Handler.
func (l SearchHandler) Handle(_ context.Context, _ Request) (Response, error) {
	return Response{}, nil
}
