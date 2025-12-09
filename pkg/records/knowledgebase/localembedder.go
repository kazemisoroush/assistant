package knowledgebase

import (
	"context"
	"math"
	"strings"
)

// LocalEmbedder is a simple embedder for POC/development
// Uses TF-IDF-like approach to generate fixed-size embeddings
type LocalEmbedder struct {
	dimensions int
	vocabulary map[string]int // Global vocabulary for consistent embeddings
}

// NewLocalEmbedder creates a new local embedder
func NewLocalEmbedder(dimensions int) Embedder {
	if dimensions <= 0 {
		dimensions = 100 // Default dimension size
	}
	return &LocalEmbedder{
		dimensions: dimensions,
		vocabulary: make(map[string]int),
	}
}

// Embed generates embeddings for text
func (le *LocalEmbedder) Embed(_ context.Context, text string) ([]float32, error) {
	terms := extractTermsForEmbedding(text)
	return le.termsToEmbedding(terms), nil
}

// EmbedBatch generates embeddings for multiple texts
func (le *LocalEmbedder) EmbedBatch(_ context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		terms := extractTermsForEmbedding(text)
		embeddings[i] = le.termsToEmbedding(terms)
	}
	return embeddings, nil
}

// Dimensions returns the dimension of the embedding vectors
func (le *LocalEmbedder) Dimensions() int {
	return le.dimensions
}

// extractTermsForEmbedding tokenizes text into terms with frequencies
func extractTermsForEmbedding(text string) map[string]float64 {
	terms := make(map[string]float64)

	// Simple tokenization: lowercase and split by whitespace/punctuation
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r < 'a' || r > 'z' && (r < '0' || r > '9')
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

// termsToEmbedding converts term frequencies to a fixed-size embedding vector
func (le *LocalEmbedder) termsToEmbedding(terms map[string]float64) []float32 {
	vector := make([]float32, le.dimensions)

	for term, freq := range terms {
		// Use hash-based indexing to map terms to vector positions
		hash := hashTerm(term)
		idx := int(hash) % le.dimensions
		vector[idx] += float32(freq)
	}

	// Normalize the vector
	magnitude := float32(0.0)
	for _, val := range vector {
		magnitude += val * val
	}
	magnitude = float32(math.Sqrt(float64(magnitude)))

	if magnitude > 0 {
		for i := range vector {
			vector[i] /= magnitude
		}
	}

	return vector
}

// hashTerm creates a simple hash for a string
func hashTerm(s string) uint32 {
	var hash uint32
	for i := 0; i < len(s); i++ {
		hash = hash*31 + uint32(s[i])
	}
	return hash
}
