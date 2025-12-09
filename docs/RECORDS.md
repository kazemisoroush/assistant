# Record Management System

A hybrid record storage and retrieval system for personal records with semantic search capabilities.

## Architecture

The system uses a **hybrid approach**:
- **Local JSON storage** for record metadata and content
- **In-memory caching** for fast access
- **Vector store interface** (ready for future AI integration with Ollama/Bedrock)
- **Flexible metadata** with type-specific structured fields

## Record Types

Currently supported record types:
- `health_visit` - Doctor visit notes
- `health_test` - Lab results and medical tests
- `health_lab` - Lab reports
- `receipt` - Purchase receipts (gas, groceries, etc.)
- `insurance` - Insurance records
- `id` - Identification records
- `travel` - Travel-related records
- `work_contract` - Employment contracts
- `tax` - Tax records
- `car` - Car-related records
- `home` - Home-related records
- `visa` - Visa and citizenship records
- `other` - Other record types

## Quick Start

### 1. Build the CLI Tool

```bash
make build
# or
go build -o bin/docs ./cmd/docs
```

### 2. Set Storage Path (Optional)

```bash
export STORAGE_PATH="./data/records"
```

### 3. Ingest Records

Ingest the sample doctor visit notes:

```bash
./bin/docs ingest health_visit ./testdata/health_visit/visit_notes.txt \
  --metadata '{"doctor_name":"Dr. Sarah Johnson","clinic":"Downtown Medical Center","visit_date":"2025-12-01T00:00:00Z"}'
```

Ingest lab results:

```bash
./bin/docs ingest health_test ./testdata/health_visit/lab_results_cbc.txt \
  --metadata '{"test_name":"Complete Blood Count","lab":"Quest Diagnostics","test_date":"2025-12-02T00:00:00Z"}'
```

Ingest chest X-ray report:

```bash
./bin/docs ingest health_test ./testdata/health_visit/chest_xray_report.txt \
  --metadata '{"test_name":"Chest X-Ray","lab":"Downtown Medical Center Imaging","test_date":"2025-12-02T00:00:00Z"}'
```

Ingest gas receipt:

```bash
./bin/docs ingest receipt ./testdata/receipts/shell_gas_receipt.txt \
  --metadata '{"vendor":"Shell","amount":60.88,"currency":"USD","category":"petrol","date":"2025-12-03T14:45:00Z"}'
```

### 4. Search Records

Search for records containing specific keywords:

```bash
./bin/docs search "cough"
./bin/docs search "Dr. Johnson"
./bin/docs search "Shell gas"
```

### 5. List Records

List all records:

```bash
./bin/docs list
```

List records by type:

```bash
./bin/docs list health_visit
./bin/docs list receipt
```

### 6. Get Record Details

Get a specific record by ID (use ID from list/search output):

```bash
./bin/docs get <record-id>
```

## Record Structure

Each record contains:
- `id` - Unique identifier (UUID)
- `type` - Record type (see list above)
- `file_path` - Original file path
- `file_name` - Original file name
- `title` - Human-readable title
- `content` - Extracted text content
- `metadata` - Type-specific structured metadata
- `tags` - Optional tags for categorization
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

### Type-Specific Metadata

#### Health Visit
```json
{
  "doctor_name": "Dr. Sarah Johnson",
  "specialty": "Internal Medicine",
  "visit_date": "2025-12-01T00:00:00Z",
  "clinic": "Downtown Medical Center",
  "diagnosis": "Post-viral syndrome",
  "prescriptions": ["Tessalon Perles 100mg"],
  "follow_up_date": "2025-12-15T00:00:00Z"
}
```

#### Receipt
```json
{
  "vendor": "Shell",
  "amount": 60.88,
  "currency": "USD",
  "date": "2025-12-03T14:45:00Z",
  "category": "petrol",
  "payment_type": "credit"
}
```

## Future Enhancements

### Phase 1: Current ✅
- Local JSON storage
- Basic keyword search
- Record ingestion CLI
- Type-specific metadata

### Phase 2: AI Integration (Next)
- Vector embeddings with Ollama/Bedrock
- Semantic search
- Natural language queries
- Record summarization

### Phase 3: Advanced Features
- PDF/DOCX/image OCR support
- Automatic metadata extraction with LLMs
- Google Drive sync
- Multi-cloud storage (S3, GCS, etc.)
- Web UI

### Phase 4: Intelligence
- Relationship mapping (link related records)
- Timeline visualization
- Automated reminders (insurance renewal, follow-ups)
- Record recommendations

## Project Structure

```
pkg/records/
├── types.go              # Record types and metadata structures
├── service.go            # Service interface
├── record_service.go   # Service implementation
├── storage/
│   └── local.go         # Local JSON storage implementation
└── vectorstore/
    └── interfaces.go    # Vector store interfaces (for future)

cmd/docs/
└── main.go              # CLI tool

testdata/
├── health_visit/        # Sample medical records
└── receipts/            # Sample receipts

data/records/          # Storage directory (gitignored)
```

## Configuration

Environment variables:
- `STORAGE_PATH` - Path to record storage directory (default: `./data/records`)

## Storage Format

Records are stored as individual JSON files in the storage directory:
```
data/records/
├── 550e8400-e29b-41d4-a716-446655440000.json
├── 550e8400-e29b-41d4-a716-446655440001.json
└── ...
```

Each file contains the complete record with all metadata and content.

## Example Queries (Future with AI)

Once vector search is integrated, you'll be able to ask:
- "When was my last doctor visit?"
- "Show me all medical tests from this year"
- "How much did I spend on gas this month?"
- "What did Dr. Johnson prescribe last time?"
- "Find all records related to my cough"

## Development

### Add a New Record Type

1. Add the type constant in `pkg/records/types.go`
2. Create metadata struct if needed
3. Update this README

### Add Vector Search

Implement one of these embedders in `pkg/records/vectorstore/`:
- `ollama_embedder.go` - For local Ollama
- `bedrock_embedder.go` - For AWS Bedrock
- `chroma_store.go` - For Chroma vector database

## License

Private project for personal use.
