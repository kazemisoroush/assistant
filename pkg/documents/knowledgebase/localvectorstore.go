package knowledgebase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/kazemisoroush/assistant/pkg/documents"
)

// LocalVectorStore is a simple in-memory vector store for POC/development
// Uses basic TF-IDF-like scoring for semantic search simulation
type LocalVectorStore struct {
	mu         sync.RWMutex
	embeddings map[string]*DocumentEmbedding // docID -> embedding
}

// DocumentEmbedding represents a document with its vector representation
type DocumentEmbedding struct {
	DocID    string
	Vector   []float64
	Terms    map[string]float64 // term -> frequency for simple vector representation
	Document *documents.Document
}

// NewLocalVectorStore creates a new local vector store instance
func NewLocalVectorStore() VectorStore {
	return &LocalVectorStore{
		embeddings: make(map[string]*DocumentEmbedding),
	}
}

// Index adds document embeddings to the vector store
// For POC, we use a simple bag-of-words approach with TF-IDF-like scoring
func (lvs *LocalVectorStore) Index(_ context.Context, doc *documents.Document) error {
	lvs.mu.Lock()
	defer lvs.mu.Unlock()

	if doc.ID == "" {
		return fmt.Errorf("document ID is required")
	}

	// Create a simple term frequency map from document content
	terms := extractTerms(doc.Content + " " + doc.Title + " " + doc.Description)

	// Create embedding
	embedding := &DocumentEmbedding{
		DocID:    doc.ID,
		Terms:    terms,
		Document: doc,
		Vector:   termsToVector(terms),
	}

	lvs.embeddings[doc.ID] = embedding
	return nil
}

// Search performs semantic similarity search using cosine similarity
func (lvs *LocalVectorStore) Search(_ context.Context, query string, limit int) ([]documents.SearchResult, error) {
	lvs.mu.RLock()
	defer lvs.mu.RUnlock()

	if len(lvs.embeddings) == 0 {
		return []documents.SearchResult{}, nil
	}

	// Create query vector
	queryTerms := extractTerms(query)
	queryVector := termsToVector(queryTerms)

	// Calculate similarity scores
	var results []documents.SearchResult
	for _, embedding := range lvs.embeddings {
		score := cosineSimilarity(queryVector, embedding.Vector)
		if score > 0 {
			results = append(results, documents.SearchResult{
				Document: *embedding.Document,
				Score:    score,
			})
		}
	}

	// Sort by score (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}

	// Apply limit
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// Delete removes document from vector store
func (lvs *LocalVectorStore) Delete(_ context.Context, docID string) error {
	lvs.mu.Lock()
	defer lvs.mu.Unlock()

	if _, exists := lvs.embeddings[docID]; !exists {
		return fmt.Errorf("document not found: %s", docID)
	}

	delete(lvs.embeddings, docID)
	return nil
}

// Close closes the vector store connection (no-op for local store)
func (lvs *LocalVectorStore) Close() error {
	return nil
}

// extractTerms tokenizes text into terms with frequencies
func extractTerms(text string) map[string]float64 {
	terms := make(map[string]float64)

	// Simple tokenization: lowercase and split by whitespace/punctuation
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})

	// Calculate term frequencies
	for _, word := range words {
		if len(word) > 2 { // Ignore very short words
			terms[word]++
		}
	}

	// Normalize frequencies
	total := float64(len(words))
	if total > 0 {
		for word := range terms {
			terms[word] = terms[word] / total
		}
	}

	return terms
}

// termsToVector converts term frequencies to a simple vector representation
func termsToVector(terms map[string]float64) []float64 {
	// For simplicity, we'll create a fixed-size vector using hash-based indexing
	vectorSize := 100
	vector := make([]float64, vectorSize)

	for term, freq := range terms {
		// Simple hash to map term to vector indices
		hash := simpleHash(term)
		idx := int(hash) % vectorSize
		vector[idx] += freq
	}

	// Normalize the vector
	magnitude := 0.0
	for _, val := range vector {
		magnitude += val * val
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude > 0 {
		for i := range vector {
			vector[i] /= magnitude
		}
	}

	return vector
}

// simpleHash creates a simple hash for a string
func simpleHash(s string) uint32 {
	var hash uint32
	for i := 0; i < len(s); i++ {
		hash = hash*31 + uint32(s[i])
	}
	return hash
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, magnitudeA, magnitudeB float64
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0.0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
