package documents

import (
	"time"
)

// DocumentType represents the category of document
type DocumentType string

const (
	DocumentTypeHealthVisit  DocumentType = "health_visit"
	DocumentTypeHealthTest   DocumentType = "health_test"
	DocumentTypeHealthLab    DocumentType = "health_lab"
	DocumentTypeReceipt      DocumentType = "receipt"
	DocumentTypeInsurance    DocumentType = "insurance"
	DocumentTypeID           DocumentType = "id"
	DocumentTypeTravel       DocumentType = "travel"
	DocumentTypeWorkContract DocumentType = "work_contract"
	DocumentTypeTax          DocumentType = "tax"
	DocumentTypeCar          DocumentType = "car"
	DocumentTypeHome         DocumentType = "home"
	DocumentTypeVisa         DocumentType = "visa"
	DocumentTypeOther        DocumentType = "other"
)

// Document represents a single document with both content and metadata
type Document struct {
	ID          string                 `json:"id"`
	Type        DocumentType           `json:"type"`
	FilePath    string                 `json:"file_path"` // Original file path
	FileName    string                 `json:"file_name"` // Original file name
	Title       string                 `json:"title"`     // Human-readable title
	Description string                 `json:"description,omitempty"`
	Content     string                 `json:"content"` // Extracted text content
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata"` // Flexible for type-specific fields
	Tags        []string               `json:"tags,omitempty"`
}

// HealthVisitMetadata contains structured data for health visits
type HealthVisitMetadata struct {
	DoctorName    string    `json:"doctor_name,omitempty"`
	Specialty     string    `json:"specialty,omitempty"`
	VisitDate     time.Time `json:"visit_date,omitempty"`
	Clinic        string    `json:"clinic,omitempty"`
	ClinicAddress string    `json:"clinic_address,omitempty"`
	Diagnosis     string    `json:"diagnosis,omitempty"`
	Symptoms      []string  `json:"symptoms,omitempty"`
	Prescriptions []string  `json:"prescriptions,omitempty"`
	FollowUpDate  time.Time `json:"follow_up_date,omitempty"`
	Notes         string    `json:"notes,omitempty"`
}

// HealthTestMetadata contains structured data for medical tests
type HealthTestMetadata struct {
	TestName  string    `json:"test_name,omitempty"`
	TestDate  time.Time `json:"test_date,omitempty"`
	OrderedBy string    `json:"ordered_by,omitempty"` // Doctor who ordered
	Lab       string    `json:"lab,omitempty"`
	Results   string    `json:"results,omitempty"`
	Notes     string    `json:"notes,omitempty"`
}

// ReceiptMetadata contains structured data for receipts
type ReceiptMetadata struct {
	Vendor      string    `json:"vendor,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Currency    string    `json:"currency,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Category    string    `json:"category,omitempty"`     // e.g., "petrol", "groceries", "medical"
	PaymentType string    `json:"payment_type,omitempty"` // e.g., "credit", "debit", "cash"
	Location    string    `json:"location,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}

// InsuranceMetadata contains structured data for insurance documents
type InsuranceMetadata struct {
	Provider     string    `json:"provider,omitempty"`
	PolicyNumber string    `json:"policy_number,omitempty"`
	Type         string    `json:"type,omitempty"` // e.g., "health", "car", "home"
	StartDate    time.Time `json:"start_date,omitempty"`
	EndDate      time.Time `json:"end_date,omitempty"`
	Premium      float64   `json:"premium,omitempty"`
	Coverage     string    `json:"coverage,omitempty"`
	ContactInfo  string    `json:"contact_info,omitempty"`
	Notes        string    `json:"notes,omitempty"`
}

// SearchResult represents a search result with relevance score
type SearchResult struct {
	Document Document `json:"document"`
	Score    float64  `json:"score"` // Relevance score (0-1)
}
