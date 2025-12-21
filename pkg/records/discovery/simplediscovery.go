// Package discovery provides implementations for discovering records based on prompts.
package discovery

import (
	"context"
	"fmt"

	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
)

// SimpleDiscovery is a basic implementation of the Discovery interface.
type SimpleDiscovery struct {
	vectorStorage knowledgebase.VectorStorage
}

// NewSimpleDiscovery creates a new instance of SimpleDiscovery.
func NewSimpleDiscovery(vectorStorage knowledgebase.VectorStorage) Discovery {
	return &SimpleDiscovery{
		vectorStorage: vectorStorage,
	}
}

// Discover implements the Discovery interface.
func (d *SimpleDiscovery) Discover(ctx context.Context, request DiscoverRequest) (DiscoverResponse, error) {
	result, err := d.vectorStorage.Search(ctx, request.Prompt, request.Limit)
	if err != nil {
		return DiscoverResponse{}, fmt.Errorf("vector storage search failed: %w", err)
	}

	hits := make([]Hit, 0, len(result))
	for _, res := range result {
		hit := Hit{
			RecordID: res.Record.ID,
			Score:    res.Score,
			Meta:     res.Record.Metadata,
			Source:   "vector",
		}
		hits = append(hits, hit)
	}

	return DiscoverResponse{
		Hits: hits,
	}, nil
}
