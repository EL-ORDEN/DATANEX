package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"datanex/internal/config"
	"datanex/internal/database"
	"datanex/internal/query"
	"datanex/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "datanex",
	Short: "Data exploration and management CLI",
	Long:  `DataNex is a modern command-line tool for exploring databases, running SQL, and managing data workflows.`,
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage database connections",
}

var queryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Execute a SQL query against a configured database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if cfg.DefaultDB == "" {
			return fmt.Errorf("no default database configured; use 'datanex db connect' first")
		}
		conn, err := findConnection(cfg, cfg.DefaultDB)
		if err != nil {
			return err
		}
		db, err := openDBConnection(conn)
		if err != nil {
			return err
		}
		defer db.Close()
		result, err := query.NewEngine(db).Execute(args[0])
		if err != nil {
			return err
		}
		printResult(result)
		return nil
	},
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open an interactive SQL shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if cfg.DefaultDB == "" {
			return fmt.Errorf("no default database configured; use 'datanex db connect' first")
		}
		conn, err := findConnection(cfg, cfg.DefaultDB)
		if err != nil {
			return err
		}
		db, err := openDBConnection(conn)
		if err != nil {
			return err
		}
		defer db.Close()
		return (&ui.SQLShell{DB: db}).Run()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd, queryCmd, shellCmd)
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func printResult(result *query.Result) {
	if result == nil {
		return
	}
	if result.HasError {
		fmt.Printf("Query failed: %s\n", result.Error)
		return
	}
	fmt.Printf("Rows: %d | Time: %s\n", result.RowCount, result.Duration)
	if len(result.Columns) == 0 {
		fmt.Println("No data returned")
		return
	}
	fmt.Println(strings.Join(result.Columns, " | "))
	for _, row := range result.Rows {
		fmt.Println(strings.Join(row, " | "))
	}
}

func openDBConnection(conn config.Connection) (*sql.DB, error) {
	switch strings.ToLower(conn.Type) {
	case "sqlite":
		db, err := database.OpenSQLite(conn.DSN)
		if err != nil {
			return nil, err
		}
		return db.DB, nil
	case "postgres", "postgresql":
		db, err := database.OpenPostgres(conn.DSN)
		if err != nil {
			return nil, err
		}
		return db.DB, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", conn.Type)
	}
}

func findConnection(cfg *config.Config, name string) (config.Connection, error) {
	for _, conn := range cfg.Connections {
		if conn.Name == name {
			return conn, nil
		}
	}
	return config.Connection{}, fmt.Errorf("connection %q not found", name)
}
