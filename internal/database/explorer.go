package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// Explorer exposes database metadata for the CLI explorer.
type Explorer struct {
	db *sql.DB
}

// NewExplorer creates a database metadata explorer.
func NewExplorer(db *sql.DB) *Explorer {
	return &Explorer{db: db}
}

// ListTables returns table names in the current schema.
func (e *Explorer) ListTables() ([]string, error) {
	if e == nil || e.db == nil {
		return nil, fmt.Errorf("explorer is not initialized")
	}
	rows, err := e.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

// DescribeTable returns column metadata for a table.
func (e *Explorer) DescribeTable(tableName string) ([]ColumnInfo, error) {
	if e == nil || e.db == nil {
		return nil, fmt.Errorf("explorer is not initialized")
	}
	if strings.TrimSpace(tableName) == "" {
		return nil, fmt.Errorf("table name is required")
	}
	query := fmt.Sprintf("PRAGMA table_info(%s)", quoteIdentifier(tableName))
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("describe table %q: %w", tableName, err)
	}
	defer rows.Close()

	var cols []ColumnInfo
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var dfltValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("scan column metadata: %w", err)
		}
		cols = append(cols, ColumnInfo{
			Name:     name,
			Type:     typ,
			NotNull:  notNull == 1,
			Default:  dfltValue.String,
			Primary:  pk == 1,
			Position: cid,
		})
	}
	return cols, rows.Err()
}

// ColumnInfo contains metadata about a database column.
type ColumnInfo struct {
	Name     string
	Type     string
	NotNull  bool
	Default  string
	Primary  bool
	Position int
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
