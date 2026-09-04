package exporter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportJSON writes a table to a JSON file as an array of objects.
func ExportJSON(db *sql.DB, tableName, filePath string) error {
	if db == nil {
		return fmt.Errorf("database handle is required")
	}
	if strings.TrimSpace(tableName) == "" {
		return fmt.Errorf("table name is required")
	}
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("file path is required")
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", quoteIdentifier(tableName)))
	if err != nil {
		return fmt.Errorf("select table rows: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("read columns: %w", err)
	}

	result := make([]map[string]any, 0)
	values := make([]any, len(cols))
	valuePtrs := make([]any, len(cols))
	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		obj := make(map[string]any, len(cols))
		for i, col := range cols {
			if values[i] == nil {
				obj[col] = nil
				continue
			}
			obj[col] = values[i]
		}
		result = append(result, obj)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate rows: %w", err)
	}

	content, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if err := os.WriteFile(filePath, content, 0o600); err != nil {
		return fmt.Errorf("write json file: %w", err)
	}
	return nil
}
