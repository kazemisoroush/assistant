// Package discovery provides implementations for discovering records based on prompts.
package discovery

import (
	"context"

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
func (sd *SimpleDiscovery) Discover(_ context.Context, _ string) error {
	// Simple discovery logic goes here.
	return nil
}
