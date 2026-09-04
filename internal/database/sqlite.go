package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteDB wraps *sql.DB and exposes a compact interface.
type SQLiteDB struct {
	*sql.DB
}

// OpenSQLite opens a SQLite database connection.
func OpenSQLite(dsn string) (*SQLiteDB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("sqlite DSN is required")
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	return &SQLiteDB{DB: db}, nil
}

// Query executes a query and returns rows.
func (d *SQLiteDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

// Exec executes a non-query statement.
func (d *SQLiteDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

// Ping validates the database connection.
func (d *SQLiteDB) Ping() error {
	return d.DB.Ping()
}

// Begin starts a transaction.
func (d *SQLiteDB) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

// HealthCheck runs a lightweight check and returns the connection time.
func (d *SQLiteDB) HealthCheck() (time.Duration, error) {
	start := time.Now()
	if err := d.Ping(); err != nil {
		return 0, fmt.Errorf("health check failed: %w", err)
	}
	return time.Since(start), nil
}
