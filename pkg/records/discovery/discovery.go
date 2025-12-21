package discovery

import "context"

// Discovery defines the interface for discovering records based on a prompt.
//
//go:generate mockgen -destination=./mocks/discovery_mock.go -package=mocks github.com/kazemisoroush/assistant/pkg/records/discovery Discovery
type Discovery interface {
	Discover(ctx context.Context, request DiscoverRequest) (DiscoverResponse, error)
}

// DiscoverRequest represents the request for a discovery operation
type DiscoverRequest struct {
	Prompt string
	Limit  int
}

// DiscoverResponse represents the response from a discovery operation
type DiscoverResponse struct {
	Hits []Hit
}

// Hit represents a single discovered record with metadata
type Hit struct {
	RecordID string
	Score    float64
	Meta     map[string]any // type/date/merchant/etc if you have it
	Source   string         // "vector", "sql", "hybrid"
}
