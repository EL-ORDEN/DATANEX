package importer

import (
	"os"
	"path/filepath"
	"testing"

	"datanex/internal/database"
)

func TestImportCSV(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	path := filepath.Join(t.TempDir(), "users.csv")
	content := "id,name,age\n1,Ada,36\n2,Linus,54\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	count, err := ImportCSV(db.DB, path, "users", true, false)
	if err != nil {
		t.Fatalf("import csv: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows imported, got %d", count)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		t.Fatalf("count users: %v", err)
	}
	defer rows.Close()
	rows.Next()
	var total int
	if err := rows.Scan(&total); err != nil {
		t.Fatalf("scan count: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total rows 2, got %d", total)
	}
}

func TestImportJSON(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	path := filepath.Join(t.TempDir(), "users.json")
	content := `[
		{"id":1,"name":"Ada","age":36},
		{"id":2,"name":"Linus","age":54}
	]`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	count, err := ImportJSON(db.DB, path, "users", true, false)
	if err != nil {
		t.Fatalf("import json: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 rows imported, got %d", count)
	}
}
