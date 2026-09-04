package analytics

import (
	"database/sql"
	"testing"

	"datanex/internal/database"
)

func TestProfileTable(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE metrics (id INTEGER PRIMARY KEY, name TEXT, score INTEGER)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := db.Exec("INSERT INTO metrics (name, score) VALUES ('Ada', 10), ('Linus', 20), ('Grace', NULL)"); err != nil {
		t.Fatalf("insert rows: %v", err)
	}

	profile, err := ProfileTable(db.DB, "metrics")
	if err != nil {
		t.Fatalf("profile table: %v", err)
	}
	if profile.RowCount != 3 {
		t.Fatalf("expected row count 3, got %d", profile.RowCount)
	}
	if profile.ColumnCount != 3 {
		t.Fatalf("expected 3 columns, got %d", profile.ColumnCount)
	}
	if len(profile.Columns) != 3 {
		t.Fatalf("expected 3 column profiles, got %d", len(profile.Columns))
	}

	var score *ColumnProfile
	for i := range profile.Columns {
		if profile.Columns[i].Name == "score" {
			score = &profile.Columns[i]
			break
		}
	}
	if score == nil {
		t.Fatal("score column profile not found")
	}
	if score.Missing != 33.33333333333333 {
		t.Fatalf("expected 33.33 missing ratio, got %v", score.Missing)
	}
	if score.Mean != "15" {
		t.Fatalf("expected mean 15, got %q", score.Mean)
	}
}

func TestProfileTableEmpty(t *testing.T) {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE empty_table (id INTEGER PRIMARY KEY, value INTEGER)"); err != nil {
		t.Fatalf("create table: %v", err)
	}

	profile, err := ProfileTable(db, "empty_table")
	if err != nil {
		t.Fatalf("profile empty table: %v", err)
	}
	if profile.RowCount != 0 {
		t.Fatalf("expected 0 rows, got %d", profile.RowCount)
	}
	if profile.ColumnCount != 2 {
		t.Fatalf("expected 2 columns, got %d", profile.ColumnCount)
	}
}
