package query

import (
	"testing"

	"datanex/internal/database"
)

func TestEngineExecutesSelect(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE metrics (id INTEGER PRIMARY KEY, value INTEGER)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec("INSERT INTO metrics (value) VALUES (10), (20), (30)"); err != nil {
		t.Fatalf("insert rows: %v", err)
	}

	e := NewEngine(db)
	result, err := e.Execute("SELECT value FROM metrics ORDER BY value")
	if err != nil {
		t.Fatalf("execute query: %v", err)
	}
	if result.RowCount != 3 {
		t.Fatalf("expected row count 3 got %d", result.RowCount)
	}
	if len(result.Columns) != 1 {
		t.Fatalf("expected 1 column got %d", len(result.Columns))
	}
}
