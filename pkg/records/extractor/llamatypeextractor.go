package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// OllamaTimeout defines the timeout for Llama API calls
const OllamaTimeout = 30 * time.Second

// LlamaTypeExtractor uses Ollama LLM to classify record types.
type LlamaTypeExtractor struct {
	ollamaURL  string
	model      string
	httpClient *http.Client
}

// NewLlamaTypeExtractor creates a new LlamaTypeExtractor instance
func NewLlamaTypeExtractor(ollamaURL, model string) TypeExtractor {
	return &LlamaTypeExtractor{
		ollamaURL: ollamaURL,
		model:     model,
		httpClient: &http.Client{
			Timeout: OllamaTimeout,
		},
	}
}

// GetType classifies the record type based on raw content
func (l *LlamaTypeExtractor) GetType(textContent string) records.RecordType {
	types := records.AllRecordTypesAsStrings()
	typesCommaSeparated := strings.Join(types, ", ")
	prompt := fmt.Sprintf("Classify the following text into exactly one of these categories: %s. Reply with ONLY the category name in lowercase. Text: %s Category:", typesCommaSeparated, textContent)

	response, err := l.callOllama(prompt)
	if err != nil {
		return records.RecordTypeOther
	}

	recordType := records.RecordType(strings.TrimSpace(strings.ToLower(response)))
	if !recordType.IsValid() {
		return records.RecordTypeOther
	}

	return recordType
}

func (l *LlamaTypeExtractor) callOllama(prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model":  l.model,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := l.httpClient.Post(
		l.ollamaURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("warning: failed to close response body: %v\n", err)
		}
	}()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	return response, nil
}
