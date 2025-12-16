// Package records provides core types and structures for record management.
package records

import (
	"time"
)

// RecordType represents the category of record
type RecordType string

// IsValid checks if the RecordType is one of the defined types
func (rt RecordType) IsValid() bool {
	validTypes := AllRecordTypes()
	for _, t := range validTypes {
		if rt == t {
			return true
		}
	}
	return false
}

// Record type constants
const (
	RecordTypeHealthVisit  RecordType = "health_visit"
	RecordTypeHealthTest   RecordType = "health_test"
	RecordTypeHealthLab    RecordType = "health_lab"
	RecordTypeReceipt      RecordType = "receipt"
	RecordTypeInsurance    RecordType = "insurance"
	RecordTypeID           RecordType = "id"
	RecordTypeTravel       RecordType = "travel"
	RecordTypeWorkContract RecordType = "work_contract"
	RecordTypeTax          RecordType = "tax"
	RecordTypeCar          RecordType = "car"
	RecordTypeHome         RecordType = "home"
	RecordTypeVisa         RecordType = "visa"
	RecordTypeOther        RecordType = "other"
)

// AllRecordTypes returns a slice of all defined record types.
func AllRecordTypes() []RecordType {
	return []RecordType{
		RecordTypeHealthVisit,
		RecordTypeHealthTest,
		RecordTypeHealthLab,
		RecordTypeReceipt,
		RecordTypeInsurance,
		RecordTypeID,
		RecordTypeTravel,
		RecordTypeWorkContract,
		RecordTypeTax,
		RecordTypeCar,
		RecordTypeHome,
		RecordTypeVisa,
		RecordTypeOther,
	}
}

// AllRecordTypesAsStrings returns all record types as string slice.
func AllRecordTypesAsStrings() []string {
	types := AllRecordTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result
}

// Record represents a single record with both content and metadata
type Record struct {
	ID        string                 `json:"id"`
	Type      RecordType             `json:"type"`
	Content   string                 `json:"content"` // Extracted text content
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata"` // Flexible for type-specific fields
	Tags      []string               `json:"tags,omitempty"`
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	Record Record  `json:"record"`
	Score  float64 `json:"score"` // Relevance score (0-1)
}
