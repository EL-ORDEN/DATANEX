package config

import "testing"

func TestConnectionRedaction(t *testing.T) {
	conn := Connection{
		Name: "users_db",
		Type: "sqlite",
		DSN:  "file:secure.db?mode=rwc&cache=shared&_auth=token",
	}

	masked := conn.RedactedDSN()
	if masked == conn.DSN {
		t.Fatal("expected DSN to be redacted")
	}
	if masked == "" {
		t.Fatal("expected masked DSN to be non-empty")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	cfg := NewConfig(tempDir)
	cfg.DefaultDB = "analytics"
	cfg.Connections = []Connection{{Name: "analytics", Type: "sqlite", DSN: "file:analytics.db"}}

	if err := cfg.Save(); err != nil {
		t.Fatalf("save config: %v", err)
	}

	loaded, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if loaded.DefaultDB != "analytics" {
		t.Fatalf("default db mismatch: got %q want %q", loaded.DefaultDB, "analytics")
	}
	if len(loaded.Connections) != 1 {
		t.Fatalf("connection count mismatch: got %d want 1", len(loaded.Connections))
	}
}
