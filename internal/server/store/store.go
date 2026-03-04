package store

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

// Store wraps the SQLite database connection.
type Store struct {
	db *sql.DB
}

// New opens a SQLite database at the given path.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite single-writer
	return &Store{db: db}, nil
}

// Init runs the initial migration.
func (s *Store) Init(migrationSQL string) error {
	slog.Info("running database migrations")
	_, err := s.db.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying *sql.DB for advanced usage.
func (s *Store) DB() *sql.DB {
	return s.db
}
