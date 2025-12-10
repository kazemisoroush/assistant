package handler

import (
	"context"
	"fmt"
)

// PrintUsageHandler handles printing CLI usage information.
type PrintUsageHandler struct{}

// NewPrintUsageHandler creates a new print usage handler.
func NewPrintUsageHandler() Handler {
	return &PrintUsageHandler{}
}

// Handle implements Handler.
func (l PrintUsageHandler) Handle(_ context.Context) {
	fmt.Println("Assistant CLI - Personal Record Management")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  assistant <command> [arguments]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  scrape              Scrape records from configured sources")
	fmt.Println("  search <query>      Search for records using a query")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  assistant scrape")
	fmt.Println("  assistant search \"health visit doctor\"")
	fmt.Println("  assistant search \"gas receipt\"")
	fmt.Println()
}
