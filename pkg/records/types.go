package records

import (
	"time"
)

// RecordType represents the category of record
type RecordType string

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

// Record represents a single record with both content and metadata
type Record struct {
	ID          string                 `json:"id"`
	Type        RecordType             `json:"type"`
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

// InsuranceMetadata contains structured data for insurance records
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
	Record Record  `json:"record"`
	Score  float64 `json:"score"` // Relevance score (0-1)
}
