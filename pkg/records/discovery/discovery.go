package discovery

import "context"

// Discovery defines the interface for discovering records based on a prompt.
//
//go:generate mockgen -destination=./mocks/discovery_mock.go -package=mocks github.com/kazemisoroush/assistant/pkg/records/discovery Discovery
type Discovery interface {
	Discover(ctx context.Context, prompt string) error
}
