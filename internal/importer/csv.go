package importer

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// ImportCSV reads a CSV file and imports it into a table.
func ImportCSV(db *sql.DB, filePath, tableName string, createTable, replace bool) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("database handle is required")
	}
	if strings.TrimSpace(filePath) == "" {
		return 0, fmt.Errorf("file path is required")
	}
	if strings.TrimSpace(tableName) == "" {
		return 0, fmt.Errorf("table name is required")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open csv file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("parse csv file: %w", err)
	}
	if len(records) < 2 {
		return 0, fmt.Errorf("csv file must include header and at least one data row")
	}

	headers := records[0]
	if createTable {
		if err := createTableFromCSV(db, tableName, headers); err != nil {
			return 0, err
		}
	}
	if replace {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", quoteIdentifier(tableName))); err != nil {
			return 0, fmt.Errorf("clear table: %w", err)
		}
	}

	cols := make([]string, 0, len(headers))
	for _, h := range headers {
		cols = append(cols, quoteIdentifier(h))
	}
	placeholders := make([]string, len(headers))
	for i := range headers {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoteIdentifier(tableName), strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	for _, row := range records[1:] {
		if len(row) != len(headers) {
			return 0, fmt.Errorf("row length mismatch in %s: expected %d, got %d", filepath.Base(filePath), len(headers), len(row))
		}
		args := make([]any, len(row))
		for i, val := range row {
			args[i] = normalizeValue(val)
		}
		if _, err := db.Exec(query, args...); err != nil {
			return 0, fmt.Errorf("insert csv row: %w", err)
		}
	}
	return len(records) - 1, nil
}

func createTableFromCSV(db *sql.DB, tableName string, headers []string) error {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", quoteIdentifier(tableName))); err != nil {
		return fmt.Errorf("drop existing table: %w", err)
	}
	parts := make([]string, 0, len(headers))
	for _, h := range headers {
		parts = append(parts, fmt.Sprintf("%s %s", quoteIdentifier(h), inferSQLiteType(h)))
	}
	query := fmt.Sprintf("CREATE TABLE %s (%s)", quoteIdentifier(tableName), strings.Join(parts, ", "))
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("create table from csv: %w", err)
	}
	return nil
}

func inferSQLiteType(colName string) string {
	_ = colName
	return "TEXT"
}

func normalizeValue(s string) any {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func isZeroValue(v any) bool {
	return reflect.ValueOf(v).IsZero()
}
