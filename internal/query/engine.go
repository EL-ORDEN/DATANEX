package query

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Result is a compact result set for CLI output.
type Result struct {
	Columns    []string
	Rows       [][]string
	RowCount   int
	Duration   time.Duration
	Query      string
	Error      string
	HasError   bool
}

// SQLQueryer exposes the minimal behavior needed by the query engine.
type SQLQueryer interface {
	Query(string, ...any) (*sql.Rows, error)
}

// Engine executes SQL queries against a database handle.
type Engine struct {
	db SQLQueryer
}

// NewEngine creates a query engine from any database implementation that supports Query.
func NewEngine(db SQLQueryer) *Engine {
	return &Engine{db: db}
}

// Execute runs a SQL query and returns a result.
func (e *Engine) Execute(query string) (*Result, error) {
	if e == nil || e.db == nil {
		return nil, fmt.Errorf("query engine is not initialized")
	}
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, fmt.Errorf("query is empty")
	}

	start := time.Now()
	rows, err := e.db.Query(trimmed)
	if err != nil {
		return &Result{Query: trimmed, Duration: time.Since(start), Error: err.Error(), HasError: true}, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return &Result{Query: trimmed, Duration: time.Since(start), Error: err.Error(), HasError: true}, err
	}

	values := make([]any, len(cols))
	valuePtrs := make([]any, len(cols))
	outputRows := make([][]string, 0)
	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return &Result{Query: trimmed, Duration: time.Since(start), Error: err.Error(), HasError: true}, err
		}
		row := make([]string, len(cols))
		for i, v := range values {
			if v == nil {
				row[i] = "NULL"
				continue
			}
			row[i] = fmt.Sprintf("%v", v)
		}
		outputRows = append(outputRows, row)
	}
	if err := rows.Err(); err != nil {
		return &Result{Query: trimmed, Duration: time.Since(start), Error: err.Error(), HasError: true}, err
	}

	result := &Result{
		Columns:  cols,
		Rows:     outputRows,
		RowCount: len(outputRows),
		Duration: time.Since(start),
		Query:    trimmed,
	}
	return result, nil
}
