// Package main is the entry point for the Assistant CLI.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kazemisoroush/assistant/pkg/config"
	"github.com/kazemisoroush/assistant/pkg/handler"
	"github.com/kazemisoroush/assistant/pkg/records/extractor"
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
	vectorStorage := knowledgebase.NewLocalVectorStorage()

	// Initialize service
	recordService := recordsvc.NewRecordService(sqliteStorage, vectorStorage)

	// Extractors
	typeExtractor := extractor.NewLlamaTypeExtractor(cfg.AI.Ollama.URL, cfg.AI.Ollama.Model)
	extractor := extractor.NewOCRContentExtractor(typeExtractor)

	// Initialize sources
	localSource := source.NewLocalSource(extractor, cfg.Sources.Local.BasePath)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	switch command {
	case "scrape":
		handler.NewLocalScraperHandler(recordService, []source.Source{localSource}).Handle(ctx)
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
