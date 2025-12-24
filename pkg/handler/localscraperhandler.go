package handler

import (
	"context"
	"fmt"

	"github.com/kazemisoroush/assistant/pkg/records/ingestor"
	"github.com/kazemisoroush/assistant/pkg/records/source"
)

var (
	// ScrapeCommandType is the command type for scraping local sources
	ScrapeCommandType = "scrape"
)

// LocalScraperHandler handles scraping records from local sources.
type LocalScraperHandler struct {
	ingestor ingestor.Ingestor
	sources  []source.Source
}

// NewLocalScraperHandler creates a new local scraper handler.
func NewLocalScraperHandler(ingestor ingestor.Ingestor, sources []source.Source) Handler {
	return &LocalScraperHandler{
		ingestor: ingestor,
		sources:  sources,
	}
}

// Handle implements Handler.
func (l LocalScraperHandler) Handle(ctx context.Context, _ Request) (Response, error) {
	recordCount := 0

	for _, src := range l.sources {
		recordChan, errChan := src.Scrape(ctx)

		for {
			select {
			case record, ok := <-recordChan:
				if !ok {
					recordChan = nil
					continue
				}
				if err := l.ingestor.Ingest(ctx, record); err != nil {
					return Response{
						Success: false,
						Errors:  []string{fmt.Sprintf("failed to ingest record from source %s: %v", src.Name(), err)},
					}, fmt.Errorf("failed to ingest record from source %s: %w", src.Name(), err)
				}
				recordCount++
			case err, ok := <-errChan:
				if !ok {
					errChan = nil
					continue
				}
				return Response{
					Success: false,
					Errors:  []string{fmt.Sprintf("error while scraping source %s: %v", src.Name(), err)},
				}, fmt.Errorf("error while scraping source %s: %w", src.Name(), err)
			}

			if recordChan == nil && errChan == nil {
				break
			}
		}
	}

	return Response{
		Success: true,
		Data: map[string]any{
			"records_ingested": recordCount,
			"sources_scraped":  len(l.sources),
		},
	}, nil
}
