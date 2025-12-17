// Package handler provides command handlers for the Assistant CLI.
package handler

import "context"

// Handler represents a command handler that processes CLI commands.
//
//go:generate mockgen -destination=./mocks/mock_handler.go -mock_names=Handler=MockHandler -package=mocks . Handler
type Handler interface {
	Handle(ctx context.Context, request Request) (Response, error)
}

// Request represents the request to a command handler.
type Request struct {
	// TODO: TBA
}

// Response represents the response from a command handler.
type Response struct {
	// TODO: TBA
}
