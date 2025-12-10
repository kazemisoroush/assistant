// Package main is the entry point for the Assistant CLI.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kazemisoroush/assistant/pkg/config"
	"github.com/kazemisoroush/assistant/pkg/handler"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	recordsvc "github.com/kazemisoroush/assistant/pkg/records/service"
	"github.com/kazemisoroush/assistant/pkg/records/source"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

func main() {
	if len(os.Args) < 2 {
		handler.NewPrintUsageHandler().Handle(context.Background())
		os.Exit(1)
	}

	command := os.Args[1]

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	localStorage, err := storage.NewLocalStorage(cfg.Records.StoragePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize vector store (using local implementation for POC)
	vectorStore := knowledgebase.NewLocalVectorStore()

	// Initialize service
	recordService := recordsvc.NewRecordService(localStorage, vectorStore)

	// Initialize sources
	sources := initializeSources(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	switch command {
	case "scrape":
		handler.NewLocalScraperHandler(recordService, sources).Handle(ctx)
	case "search":
		handler.NewSearchHandler(recordService).Handle(ctx)
	case "help", "-h", "--help":
		handler.NewPrintUsageHandler().Handle(ctx)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		handler.NewPrintUsageHandler().Handle(ctx)
		os.Exit(1)
	}
}

func initializeSources(cfg config.Config) []source.Source {
	var sources []source.Source

	localSource := source.NewLocalSource(
		"local",
		cfg.Records.Sources.Local.BasePath,
		cfg.Records.Sources.Local.Enabled,
	)
	sources = append(sources, localSource)

	return sources
}
