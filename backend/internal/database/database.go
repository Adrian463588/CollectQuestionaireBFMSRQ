package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend/internal/config"

	_ "github.com/lib/pq"
)

// DB holds the database connection.
type DB struct {
	Connection *sql.DB
}

// NewDB creates a new database connection.
func NewDB(cfg *config.Config) (*DB, error) {
	dbURL := cfg.GetDatabaseURL()

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return &DB{Connection: conn}, nil
}

// Close closes the database connection.
func (db *DB) Close() {
	if err := db.Connection.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// Migrate runs all database migrations idempotently.
// Sprint 2 adds soft-delete columns (deleted_at) to all three tables.
// ALTER TABLE ... ADD COLUMN IF NOT EXISTS is idempotent on PostgreSQL ≥9.6.
func (db *DB) Migrate() error {
	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "create_participants",
			sql: `CREATE TABLE IF NOT EXISTS participants (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name VARCHAR(255) NOT NULL,
				age INTEGER NOT NULL CHECK (age >= 15),
				gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE NULL
			)`,
		},
		{
			name: "create_responses",
			sql: `CREATE TABLE IF NOT EXISTS responses (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				participant_id UUID REFERENCES participants(id) ON DELETE CASCADE,
				questionnaire_type VARCHAR(50) NOT NULL CHECK (questionnaire_type IN ('srq29', 'ipip-bfm-50')),
				answers JSONB NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE NULL
			)`,
		},
		{
			name: "create_scores",
			sql: `CREATE TABLE IF NOT EXISTS scores (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				participant_id UUID REFERENCES participants(id) ON DELETE CASCADE,
				srq_score JSONB,
				ipip_score JSONB,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP WITH TIME ZONE NULL
			)`,
		},
		// Sprint 2: add soft-delete column to pre-existing tables (idempotent)
		{
			name: "add_deleted_at_participants",
			sql:  `ALTER TABLE participants ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE NULL`,
		},
		{
			name: "add_deleted_at_responses",
			sql:  `ALTER TABLE responses ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE NULL`,
		},
		{
			name: "add_deleted_at_scores",
			sql:  `ALTER TABLE scores ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE NULL`,
		},
	}

	for _, m := range migrations {
		if _, err := db.Connection.Exec(m.sql); err != nil {
			return fmt.Errorf("migration %q failed: %w", m.name, err)
		}
		log.Printf("Migration %q applied", m.name)
	}

	log.Println("All database migrations completed successfully")
	return nil
}
