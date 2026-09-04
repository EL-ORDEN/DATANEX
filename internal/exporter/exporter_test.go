package exporter

import (
	"os"
	"path/filepath"
	"testing"

	"datanex/internal/database"
)

func TestExportCSVAndJSON(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec("INSERT INTO users (name, age) VALUES ('Ada', 36), ('Linus', 54)"); err != nil {
		t.Fatalf("insert rows: %v", err)
	}

	csvPath := filepath.Join(t.TempDir(), "users.csv")
	if err := ExportCSV(db.DB, "users", csvPath); err != nil {
		t.Fatalf("export csv: %v", err)
	}
	if _, err := os.Stat(csvPath); err != nil {
		t.Fatalf("csv file missing: %v", err)
	}

	jsonPath := filepath.Join(t.TempDir(), "users.json")
	if err := ExportJSON(db.DB, "users", jsonPath); err != nil {
		t.Fatalf("export json: %v", err)
	}
	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("json file missing: %v", err)
	}
}
