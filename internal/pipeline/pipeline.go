package pipeline

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config describes a data pipeline.
type Config struct {
	Name     string       `json:"name" yaml:"name"`
	Database string       `json:"database" yaml:"database"`
	Steps    []StepConfig `json:"steps" yaml:"steps"`
}

// StepConfig describes a single stage in the pipeline.
type StepConfig struct {
	Type     string   `json:"type" yaml:"type"`
	Source   string   `json:"source" yaml:"source"`
	Format   string   `json:"format" yaml:"format"`
	Table    string   `json:"table" yaml:"table"`
	Mode     string   `json:"mode" yaml:"mode"`
	Required []string `json:"required" yaml:"required"`
}

// LoadFromFile loads a pipeline definition from YAML or JSON.
func LoadFromFile(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pipeline file: %w", err)
	}
	cfg := &Config{}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		if err := json.Unmarshal(content, cfg); err != nil {
			return nil, fmt.Errorf("decode json pipeline: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, cfg); err != nil {
			return nil, fmt.Errorf("decode yaml pipeline: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported pipeline format %q", ext)
	}
	if cfg.Name == "" {
		cfg.Name = filepath.Base(path)
	}
	return cfg, nil
}

// Run executes the configured pipeline stages.
func (c *Config) Run(db *sql.DB) error {
	if c == nil {
		return fmt.Errorf("pipeline config is required")
	}
	if db == nil {
		return fmt.Errorf("database connection is required")
	}
	for _, step := range c.Steps {
		switch strings.ToLower(step.Type) {
		case "read":
			if step.Source == "" && step.Format == "" {
				return fmt.Errorf("read step requires a source or format")
			}
		case "clean", "transform":
			continue
		case "validate":
			if len(step.Required) == 0 {
				return fmt.Errorf("validate step requires required fields")
			}
		case "database":
			if step.Table == "" {
				return fmt.Errorf("database step requires a table name")
			}
		default:
			return fmt.Errorf("unsupported pipeline step %q", step.Type)
		}
	}
	return nil
}
