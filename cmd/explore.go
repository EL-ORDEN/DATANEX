package main

import (
	"database/sql"
	"fmt"

	"datanex/internal/config"
	"datanex/internal/database"

	"github.com/spf13/cobra"
)

var exploreCmd = &cobra.Command{
	Use:   "explore [table]",
	Short: "Inspect a table and its schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if cfg.DefaultDB == "" {
			return fmt.Errorf("no default database configured")
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
		return printTableMetadata(db, args[0])
	},
}

func printTableMetadata(db *sql.DB, tableName string) error {
	explorer := database.NewExplorer(db)
	cols, err := explorer.DescribeTable(tableName)
	if err != nil {
		return err
	}
	if len(cols) == 0 {
		fmt.Printf("Table %q does not exist or has no columns.\n", tableName)
		return nil
	}
	fmt.Printf("Table: %s\n", tableName)
	fmt.Println("Column         Type            Not Null   PK")
	fmt.Println("--------------------------------------------------")
	for _, col := range cols {
		fmt.Printf("%-14s %-15s %-9t %-4t\n", col.Name, col.Type, col.NotNull, col.Primary)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(exploreCmd)
}
