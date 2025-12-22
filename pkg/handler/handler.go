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
	// Command is the name of the command being executed
	Command string

	// Data contains the parsed/typed payload for the command
	// Handlers can type-assert this to their specific input types
	Data any
}

// Response represents the response from a command handler.
type Response struct {
	// Success indicates whether the command executed successfully
	Success bool

	// Data contains the result payload from the handler
	// Can be marshaled to JSON or formatted for CLI output
	Data any

	// Errors contains any error details (can be multiple validation errors)
	Errors []string
}
