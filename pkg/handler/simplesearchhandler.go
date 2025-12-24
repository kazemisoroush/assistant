package handler

import (
	"context"
	"fmt"

	"github.com/kazemisoroush/assistant/pkg/records/discovery"
)

const (
	// DefaultSearchLimit is the default maximum number of search results
	DefaultSearchLimit = 10

	// SimpleSearchCommandType is the command type for simple search operations
	SimpleSearchCommandType = "search"
)

// SimpleSearchHandler handles searching for records.
type SimpleSearchHandler struct {
	discovery discovery.Discovery
}

// NewSimpleSearchHandler creates a new search handler.
func NewSimpleSearchHandler(discovery discovery.Discovery) Handler {
	return &SimpleSearchHandler{
		discovery: discovery,
	}
}

// Handle implements Handler for search operations.
func (h *SimpleSearchHandler) Handle(ctx context.Context, request Request) (Response, error) {
	// Extract search prompt from request data
	prompt, ok := request.Data.(string)
	if !ok || prompt == "" {
		return Response{
			Success: false,
			Errors:  []string{"search prompt is required"},
		}, fmt.Errorf("search prompt is required")
	}

	// Perform discovery with default limit
	discoverRequest := discovery.DiscoverRequest{
		Prompt: prompt,
		Limit:  DefaultSearchLimit,
	}

	discoverResponse, err := h.discovery.Discover(ctx, discoverRequest)
	if err != nil {
		return Response{
			Success: false,
			Errors:  []string{fmt.Sprintf("search failed: %v", err)},
		}, fmt.Errorf("search failed: %w", err)
	}

	// Return successful response with hits
	return Response{
		Success: true,
		Data:    discoverResponse.Hits,
	}, nil
}
