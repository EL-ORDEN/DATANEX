package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"datanex/internal/database"
)

func TestLoadPipelineConfigYAML(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "users.yaml")
	content := `name: users_pipeline
database: demo
steps:
  - type: read
    source: ./data/users.csv
    format: csv
  - type: validate
    required: [id, name]
  - type: database
    table: users
    mode: replace
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write pipeline file: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("load pipeline: %v", err)
	}
	if cfg.Name != "users_pipeline" {
		t.Fatalf("expected name users_pipeline, got %q", cfg.Name)
	}
	if len(cfg.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(cfg.Steps))
	}
}

func TestPipelineRun(t *testing.T) {
	db, err := database.OpenSQLite("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	cfg := &Config{
		Name:     "demo",
		Database: "demo",
		Steps: []StepConfig{
			{Type: "read", Source: "", Format: "csv"},
			{Type: "clean"},
			{Type: "validate", Required: []string{"id", "name"}},
			{Type: "database", Table: "users", Mode: "replace"},
		},
	}

	if err := cfg.Run(db.DB); err != nil {
		t.Fatalf("run pipeline: %v", err)
	}
}
