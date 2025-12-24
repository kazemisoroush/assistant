// Package main is the entry point for the Assistant CLI.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kazemisoroush/assistant/pkg/config"
	"github.com/kazemisoroush/assistant/pkg/handler"
	"github.com/kazemisoroush/assistant/pkg/records/discovery"
	"github.com/kazemisoroush/assistant/pkg/records/extractor"
	"github.com/kazemisoroush/assistant/pkg/records/ingestor"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	"github.com/kazemisoroush/assistant/pkg/records/source"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize storage
	sqliteStorage, err := storage.NewSQLiteStorage(cfg.SQLitePath)
	if err != nil {
		slog.Error("Failed to initialize local storage", "error", err)
		os.Exit(1)
	}

	// Initialize vector store (using local implementation for POC)
	localVectorStorage := knowledgebase.NewLocalVectorStorage()

	// Initialize service
	recordService := ingestor.NewRecordIngestor(sqliteStorage, localVectorStorage)

	// Extractors
	typeExtractor := extractor.NewLlamaTypeExtractor(cfg.AI.Ollama.URL, cfg.AI.Ollama.Model)
	extractor := extractor.NewOCRContentExtractor(typeExtractor)

	// Initialize sources
	localSource := source.NewLocalSource(extractor, cfg.Sources.Local.BasePath)

	// Initialize discovery service
	discoveryService := discovery.NewSimpleDiscovery(localVectorStorage)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	switch command {
	case handler.ScrapeCommandType:
		hand := handler.NewLocalScraperHandler(recordService, []source.Source{localSource})
		resp, err := hand.Handle(ctx, handler.Request{
			Command: handler.ScrapeCommandType,
		})
		if err != nil {
			slog.Error("Scrape command failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Scrape command completed", "response", resp)
	case handler.SimpleSearchCommandType:
		hand := handler.NewSimpleSearchHandler(discoveryService)
		resp, err := hand.Handle(ctx, handler.Request{
			Command: handler.SimpleSearchCommandType,
			Data:    os.Args[2],
		})
		if err != nil {
			slog.Error("Search command failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Search command completed", "response", resp)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		os.Exit(1)
	}
}
