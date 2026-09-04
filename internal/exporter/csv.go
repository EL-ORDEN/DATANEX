package exporter

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// ExportCSV writes a table to a CSV file.
func ExportCSV(db *sql.DB, tableName, filePath string) error {
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

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create csv file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	if err := writer.Write(cols); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}

	values := make([]any, len(cols))
	valuePtrs := make([]any, len(cols))
	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("scan row: %w", err)
		}
		record := make([]string, len(cols))
		for i, v := range values {
			if v == nil {
				record[i] = ""
				continue
			}
			record[i] = fmt.Sprintf("%v", v)
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush csv file: %w", err)
	}
	return rows.Err()
}

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
