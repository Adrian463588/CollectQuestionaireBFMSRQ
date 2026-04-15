package database

import (
	"database/sql"
	"fmt"
	"log"

	"backend/internal/config"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	Connection *sql.DB
}

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*DB, error) {
	dbURL := cfg.GetDatabaseURL()

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	return &DB{Connection: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if err := db.Connection.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS participants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			age INTEGER NOT NULL CHECK (age >= 15),
			gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS responses (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			participant_id UUID REFERENCES participants(id) ON DELETE CASCADE,
			questionnaire_type VARCHAR(50) NOT NULL CHECK (questionnaire_type IN ('srq29', 'ipip-bfm-50')),
			answers JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS scores (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			participant_id UUID REFERENCES participants(id) ON DELETE CASCADE,
			srq_score JSONB,
			ipip_score JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, migration := range migrations {
		if _, err := db.Connection.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w\nSQL: %s", err, migration)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
