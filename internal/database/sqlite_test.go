package database

import "testing"

func TestSQLiteOpenAndQuery(t *testing.T) {
	conn, err := OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)"); err != nil {
		t.Fatalf("create table: %v", err)
	}
	if _, err := conn.Exec("INSERT INTO users (name) VALUES ('Ada'), ('Linus')"); err != nil {
		t.Fatalf("insert rows: %v", err)
	}

	rows, err := conn.Query("SELECT id, name FROM users ORDER BY id")
	if err != nil {
		t.Fatalf("query rows: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 rows got %d", count)
	}
}
