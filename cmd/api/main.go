// Package main is the entry point for the Assistant API server.
package main

import (
	"fmt"
	"log"

	"github.com/kazemisoroush/assistant/pkg/config"
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
	fmt.Println("Hello from Assistant API!")
}
