# Personal AI Assistant

**Instant access to your life's data, intelligently organized.**

## Purpose

A personal AI assistant that connects to your personal data stores to provide instant, natural-language access to all your important documents—medical records, receipts, IDs, travel documents, and family information. Instead of navigating folders for minutes, ask a question and get answers in seconds.

## Features

### ✅ Current Features (v0.1)

- **Document Storage**: Local JSON-based storage with in-memory caching
- **Flexible Metadata**: Type-specific structured metadata for better organization
- **Keyword Search**: Basic text search across all document content
- **CLI Tool**: Simple command-line interface for document management
- **Multiple Document Types**: Health visits, medical tests, receipts, insurance, IDs, and more

### 🚧 Coming Soon

- **Semantic Search**: Vector embeddings with Ollama/AWS Bedrock
- **Natural Language Queries**: Ask questions in plain English
- **OCR Support**: Extract text from PDFs, images, and scanned documents
- **Google Drive Sync**: Automatically sync from your Google Drive
- **Web UI**: User-friendly web interface
- **Smart Reminders**: Get notified about insurance renewals, follow-ups, etc.

## Quick Start

### 1. Build the Application

```bash
make build
```

### 2. Run the Demo

```bash
./demo.sh
```

This will ingest sample documents (doctor visit notes, lab results, and a gas receipt) and demonstrate search capabilities.

### 3. Try It Yourself

```bash
# Ingest a document
./bin/docs ingest receipt /path/to/receipt.txt \
  --metadata '{"vendor":"Store Name","amount":50.00,"currency":"USD","category":"groceries","date":"2025-12-05T12:00:00Z"}'

# Search for documents
./bin/docs search "keyword"

# List all documents
./bin/docs list

# List by type
./bin/docs list health_visit

# Get document details
./bin/docs get <document-id>
```

## Documentation

- **[Document Management Guide](docs/DOCUMENTS.md)** - Complete documentation for the document system
- Architecture overview, document types, metadata structures, and future roadmap

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build binaries
make build

# Complete CI pipeline
make ci
```

## Status

🚧 **Early Development** - Core document storage and retrieval working. AI features coming next.

