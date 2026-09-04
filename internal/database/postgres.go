package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresDB wraps *sql.DB and exposes a compact interface.
type PostgresDB struct {
	*sql.DB
}

// OpenPostgres opens a PostgreSQL database connection using the pgx stdlib driver.
func OpenPostgres(dsn string) (*PostgresDB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres DSN is required")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres database: %w", err)
	}
	return &PostgresDB{DB: db}, nil
}

// Query executes a query.
func (d *PostgresDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

// Exec executes a statement.
func (d *PostgresDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

// HealthCheck ensures the database is reachable and reports latency.
func (d *PostgresDB) HealthCheck() (time.Duration, error) {
	start := time.Now()
	if err := d.Ping(); err != nil {
		return 0, fmt.Errorf("health check failed: %w", err)
	}
	return time.Since(start), nil
}
