// Package storage provides persistent storage implementations for records.
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	// Import sqlite3 driver for database/sql
	_ "github.com/mattn/go-sqlite3"

	"github.com/kazemisoroush/assistant/pkg/records"
)

// SQLiteStorage implements the storage.SQLiteStorage interface using SQLite
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance with the given database path.
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	s := &SQLiteStorage{db: db}

	// Initialize schema
	if err := s.initSchema(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

// initSchema creates the necessary tables
func (s SQLiteStorage) initSchema() error {
	schema := `
    CREATE TABLE IF NOT EXISTS records (
        id TEXT PRIMARY KEY,
        type TEXT NOT NULL,
        content TEXT NOT NULL,
        metadata TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_records_type ON records(type);
    CREATE INDEX IF NOT EXISTS idx_records_created_at ON records(created_at);
    `

	_, err := s.db.Exec(schema)
	return err
}

// Store saves a record
func (s SQLiteStorage) Store(ctx context.Context, rec records.Record) error {
	metadata, err := json.Marshal(rec.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
        INSERT INTO records (id, type, content, metadata, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	_, err = s.db.ExecContext(ctx, query,
		rec.ID,
		rec.Type,
		rec.Content,
		string(metadata),
		rec.CreatedAt,
		rec.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to store record: %w", err)
	}

	return nil
}

// Get retrieves a record by ID
func (s SQLiteStorage) Get(ctx context.Context, id string) (records.Record, error) {
	query := `
        SELECT id, type, content, metadata, created_at, updated_at
        FROM records
        WHERE id = ?
    `

	var rec records.Record
	var metadataJSON string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&rec.ID,
		&rec.Type,
		&rec.Content,
		&metadataJSON,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return records.Record{}, fmt.Errorf("record not found: %s", id)
	}
	if err != nil {
		return records.Record{}, fmt.Errorf("failed to get record: %w", err)
	}

	if err := json.Unmarshal([]byte(metadataJSON), &rec.Metadata); err != nil {
		return records.Record{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return rec, nil
}

// List returns all records with optional type filter
func (s SQLiteStorage) List(ctx context.Context, recType records.RecordType) ([]records.Record, error) {
	var query string
	var args []interface{}

	if recType != "" {
		query = `
            SELECT id, type, content, metadata, created_at, updated_at
            FROM records
            WHERE type = ?
            ORDER BY created_at DESC
        `
		args = append(args, recType)
	} else {
		query = `
            SELECT id, type, content, metadata, created_at, updated_at
            FROM records
            ORDER BY created_at DESC
        `
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("warning: failed to close rows: %v\n", err)
		}
	}()

	var recs []records.Record
	for rows.Next() {
		var rec records.Record
		var metadataJSON string

		if err := rows.Scan(
			&rec.ID,
			&rec.Type,
			&rec.Content,
			&metadataJSON,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}

		if err := json.Unmarshal([]byte(metadataJSON), &rec.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		recs = append(recs, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating records: %w", err)
	}

	return recs, nil
}

// Update updates an existing record
func (s SQLiteStorage) Update(ctx context.Context, rec records.Record) error {
	metadata, err := json.Marshal(rec.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
        UPDATE records
        SET type = ?, content = ?, metadata = ?, updated_at = ?
        WHERE id = ?
    `

	result, err := s.db.ExecContext(ctx, query,
		rec.Type,
		rec.Content,
		string(metadata),
		rec.UpdatedAt,
		rec.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("record not found: %s", rec.ID)
	}

	return nil
}

// Delete removes a record
func (s SQLiteStorage) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM records WHERE id = ?`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("record not found: %s", id)
	}

	return nil
}

// Close closes the database connection
func (s SQLiteStorage) Close() error {
	return s.db.Close()
}
