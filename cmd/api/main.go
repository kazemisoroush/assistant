// Package main is the entry point for the Assistant API server.
package main

import (
	"fmt"
	"log"

	"github.com/kazemisoroush/assistant/pkg/config"
	"github.com/kazemisoroush/assistant/pkg/records/knowledgebase"
	recordsvc "github.com/kazemisoroush/assistant/pkg/records/service"
	"github.com/kazemisoroush/assistant/pkg/records/storage"
)

// @title Assistant API
// @version 1.0
// @description API for the Assistant application
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Assistant API starting with log level: %s\n", cfg.LogLevel)

	// Initialize storage
	localStorage, err := storage.NewSQLiteStorage(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize vector store (using local implementation for POC)
	vectorStorage := knowledgebase.NewLocalVectorStorage()

	// Initialize service
	recordService := recordsvc.NewRecordService(localStorage, vectorStorage)

	// TODO: Setup HTTP server and routes using the service or handlers
	// For now, just verify initialization
	_ = recordService

	fmt.Println("Assistant API initialized successfully!")
	fmt.Println("Service ready for API endpoints")
}
