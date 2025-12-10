package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/kazemisoroush/assistant/pkg/records/source"
)

// LocalScraperHandler handles scraping records from local sources.
type LocalScraperHandler struct {
	service records.Service
	sources []source.Source
}

// NewLocalScraperHandler creates a new local scraper handler.
func NewLocalScraperHandler(service records.Service, sources []source.Source) Handler {
	return &LocalScraperHandler{
		service: service,
		sources: sources,
	}
}

// Handle implements Handler.
func (l LocalScraperHandler) Handle(ctx context.Context) {
	fmt.Println("🔍 Starting scrape operation...")
	fmt.Println()

	// Scrape from all enabled sources
	var totalScraped, totalFailed int

	for _, src := range l.sources {
		if !src.IsEnabled() {
			fmt.Printf("⏭️  Skipping disabled source: %s\n", src.Name())
			continue
		}

		fmt.Printf("📦 Source: %s\n", src.Name())

		recordChan, errChan := src.Scrape(ctx)
		sourceScraped := 0
		sourceFailed := 0

		// Process records and errors
		for {
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stderr, "Scrape cancelled: %v\n", ctx.Err())
				os.Exit(1)

			case rec, ok := <-recordChan:
				if !ok {
					// Channel closed, no more records
					recordChan = nil
					if errChan == nil {
						// Both channels closed
						goto sourceDone
					}
					continue
				}

				// Ingest the record
				if err := l.service.Ingest(ctx, rec); err != nil {
					fmt.Fprintf(os.Stderr, "   ❌ Failed to ingest %s: %v\n", rec.Title, err)
					sourceFailed++
				} else {
					sourceScraped++
				}

			case err, ok := <-errChan:
				if !ok {
					// Error channel closed
					errChan = nil
					if recordChan == nil {
						// Both channels closed
						goto sourceDone
					}
					continue
				}

				fmt.Fprintf(os.Stderr, "   ⚠️  Scrape error: %v\n", err)
			}
		}

	sourceDone:
		fmt.Printf("   ✅ Scraped: %d records\n", sourceScraped)
		if sourceFailed > 0 {
			fmt.Printf("   ❌ Failed: %d records\n", sourceFailed)
		}
		fmt.Println()

		totalScraped += sourceScraped
		totalFailed += sourceFailed
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📊 Total: %d records scraped, %d failed\n", totalScraped, totalFailed)

	if totalFailed > 0 {
		os.Exit(1)
	}
}
