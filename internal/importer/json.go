package importer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportJSON reads a JSON array of objects and imports it into a table.
func ImportJSON(db *sql.DB, filePath, tableName string, createTable, replace bool) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database handle is required")
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("read json file: %w", err)
	}
	var rows []map[string]any
	if err := json.Unmarshal(content, &rows); err != nil {
		return 0, fmt.Errorf("decode json file: %w", err)
	}
	if len(rows) == 0 {
		return 0, nil
	}

	keys := make([]string, 0, len(rows[0]))
	for key := range rows[0] {
		keys = append(keys, key)
	}
	if createTable {
		if err := createTableFromJSON(db, tableName, keys); err != nil {
			return 0, err
		}
	}
	if replace {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", quoteIdentifier(tableName))); err != nil {
			return 0, fmt.Errorf("clear table: %w", err)
		}
	}

	placeholders := make([]string, len(keys))
	for i := range keys {
		placeholders[i] = "?"
	}
	quotedKeys := make([]string, len(keys))
	for i, key := range keys {
		quotedKeys[i] = quoteIdentifier(key)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoteIdentifier(tableName), strings.Join(quotedKeys, ", "), strings.Join(placeholders, ", "))
	for _, row := range rows {
		values := make([]any, len(keys))
		for i, key := range keys {
			values[i] = normalizeValue(fmt.Sprintf("%v", row[key]))
		}
		if _, err := db.Exec(query, values...); err != nil {
			return 0, fmt.Errorf("insert json row: %w", err)
		}
	}
	return len(rows), nil
}

func createTableFromJSON(db *sql.DB, tableName string, columns []string) error {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", quoteIdentifier(tableName))); err != nil {
		return fmt.Errorf("drop existing table: %w", err)
	}
	parts := make([]string, 0, len(columns))
	for _, col := range columns {
		parts = append(parts, fmt.Sprintf("%s TEXT", quoteIdentifier(col)))
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s)", quoteIdentifier(tableName), strings.Join(parts, ", "))
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create table from json: %w", err)
	}
	return nil
}
