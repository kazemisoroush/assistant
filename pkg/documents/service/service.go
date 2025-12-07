package documents

import (
	"context"
	"fmt"
	"strings"

	"github.com/kazemisoroush/assistant/pkg/documents"
	"github.com/kazemisoroush/assistant/pkg/documents/storage"
)

// DocumentService implements the Service interface
type DocumentService struct {
	storage storage.Storage
	// vectorStore will be added later for semantic search
	// vectorStore VectorStore
}

// NewDocumentService creates a new document service
func NewDocumentService(storage storage.Storage) documents.Service {
	return &DocumentService{
		storage: storage,
	}
}

// Ingest processes and stores a document
func (s *DocumentService) Ingest(ctx context.Context, doc *documents.Document) error {
	// Validate document
	if doc.ID == "" {
		return fmt.Errorf("document ID is required")
	}
	if doc.Type == "" {
		return fmt.Errorf("document type is required")
	}
	if doc.Content == "" && doc.FilePath == "" {
		return fmt.Errorf("document must have either content or file path")
	}

	// Initialize metadata map if nil
	if doc.Metadata == nil {
		doc.Metadata = make(map[string]interface{})
	}

	// Store the document
	if err := s.storage.Store(ctx, doc); err != nil {
		return fmt.Errorf("failed to store document: %w", err)
	}

	// TODO: Index in vector store for semantic search
	// if s.vectorStore != nil {
	//     if err := s.vectorStore.Index(ctx, doc); err != nil {
	//         return fmt.Errorf("failed to index document: %w", err)
	//     }
	// }

	return nil
}

// Search performs search with optional metadata filters
func (s *DocumentService) Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]documents.SearchResult, error) {
	// For now, use basic keyword search from storage
	// Later this will use vector store for semantic search
	if localStorage, ok := s.storage.(interface {
		Search(ctx context.Context, query string, filters map[string]interface{}, limit int) ([]documents.SearchResult, error)
	}); ok {
		return localStorage.Search(ctx, query, filters, limit)
	}

	return nil, fmt.Errorf("search not supported by current storage implementation")
}

// GetByID retrieves a document by its ID
func (s *DocumentService) GetByID(ctx context.Context, id string) (*documents.Document, error) {
	return s.storage.Get(ctx, id)
}

// List returns all documents with optional type filter
func (s *DocumentService) List(ctx context.Context, docType documents.DocumentType) ([]*documents.Document, error) {
	return s.storage.List(ctx, docType)
}

// Update updates an existing document
func (s *DocumentService) Update(ctx context.Context, doc *documents.Document) error {
	if doc.ID == "" {
		return fmt.Errorf("document ID is required")
	}

	if err := s.storage.Update(ctx, doc); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	// TODO: Update in vector store
	// if s.vectorStore != nil {
	//     if err := s.vectorStore.Index(ctx, doc); err != nil {
	//         return fmt.Errorf("failed to reindex document: %w", err)
	//     }
	// }

	return nil
}

// Delete removes a document
func (s *DocumentService) Delete(ctx context.Context, id string) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// TODO: Delete from vector store
	// if s.vectorStore != nil {
	//     if err := s.vectorStore.Delete(ctx, id); err != nil {
	//         return fmt.Errorf("failed to delete from vector store: %w", err)
	//     }
	// }

	return nil
}

// ExtractTextFromFile is a helper function to extract text content from various file types
// For now, it just reads plain text. Later we can add PDF, DOCX, image OCR support
func ExtractTextFromFile(_ string) (string, error) {
	// TODO: Implement based on file type
	// - .txt: read directly
	// - .pdf: use pdf library
	// - .docx: use docx library
	// - .jpg, .png: use OCR
	return "", fmt.Errorf("not implemented yet")
}

// NormalizeContent performs basic text normalization
func NormalizeContent(content string) string {
	// Trim whitespace
	content = strings.TrimSpace(content)
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	return content
}
