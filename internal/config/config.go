package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Connection stores a database connection definition in the local config.
type Connection struct {
	Name string `json:"name"`
	Type string `json:"type"`
	DSN  string `json:"dsn"`
}

// RedactedDSN hides credentials while preserving enough structure for diagnostics.
func (c Connection) RedactedDSN() string {
	if c.DSN == "" {
		return ""
	}
	parts := strings.Split(c.DSN, "@")
	if len(parts) == 2 {
		prefix := parts[0]
		if idx := strings.LastIndex(prefix, ":"); idx != -1 {
			prefix = prefix[:idx+1] + "***"
		}
		return prefix + "@" + parts[1]
	}
	if strings.Contains(c.DSN, "password=") || strings.Contains(c.DSN, "pass=") {
		return redactKeyValuePairs(c.DSN)
	}
	if strings.Contains(c.DSN, "://") {
		return redactURL(c.DSN)
	}
	return "***redacted***"
}

func redactKeyValuePairs(raw string) string {
	segments := strings.Split(raw, "&")
	for i, seg := range segments {
		if strings.Contains(seg, "password=") || strings.Contains(seg, "pass=") {
			kv := strings.SplitN(seg, "=", 2)
			if len(kv) == 2 {
				segments[i] = kv[0] + "=***"
			}
		}
	}
	return strings.Join(segments, "&")
}

func redactURL(raw string) string {
	if !strings.Contains(raw, "://") {
		return raw
	}
	prefix, rest, found := strings.Cut(raw, "://")
	if !found {
		return raw
	}
	if strings.Contains(rest, "/") {
		hostAndPath := strings.SplitN(rest, "/", 2)
		if len(hostAndPath) == 2 {
			return prefix + "://" + hostAndPath[0] + "/" + hostAndPath[1]
		}
	}
	return prefix + "://***redacted***"
}

// Config stores application settings, saved in a JSON file.
type Config struct {
	BaseDir    string       `json:"-"`
	Connections []Connection `json:"connections"`
	DefaultDB   string       `json:"default_db"`
	LogLevel    string       `json:"log_level"`
	OutputDir   string       `json:"output_dir"`
}

// NewConfig creates a config with defaults.
func NewConfig(baseDir string) *Config {
	if baseDir == "" {
		baseDir = defaultConfigDir()
	}
	return &Config{
		BaseDir:     baseDir,
		Connections: []Connection{},
		DefaultDB:   "",
		LogLevel:    "info",
		OutputDir:   filepath.Join(baseDir, "output"),
	}
}

// FilePath returns the config file path.
func (c *Config) FilePath() string {
	baseDir := c.BaseDir
	if baseDir == "" {
		baseDir = defaultConfigDir()
	}
	return filepath.Join(baseDir, "config.json")
}

// Save persists the configuration to disk.
func (c *Config) Save() error {
	baseDir := c.BaseDir
	if baseDir == "" {
		baseDir = defaultConfigDir()
	}
	path := filepath.Join(baseDir, "config.json")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, content, 0o600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

// LoadConfig reads config from the standard config path.
func LoadConfig(baseDir string) (*Config, error) {
	if baseDir == "" {
		baseDir = defaultConfigDir()
	}
	path := filepath.Join(baseDir, "config.json")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := NewConfig(baseDir)
			return cfg, nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	cfg := NewConfig(baseDir)
	if err := json.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	cfg.BaseDir = baseDir
	return cfg, nil
}

func defaultConfigDir() string {
	if dir := os.Getenv("DATANEX_CONFIG_DIR"); dir != "" {
		return dir
	}
	switch {
	case os.Getenv("APPDATA") != "":
		return filepath.Join(os.Getenv("APPDATA"), "DataNex")
	case os.Getenv("HOME") != "":
		return filepath.Join(os.Getenv("HOME"), ".datanex")
	default:
		return filepath.Join(".", ".datanex")
	}
}
