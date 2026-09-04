package main

import (
	"fmt"
	"os"
	"strings"

	"datanex/internal/analytics"
	"datanex/internal/config"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [table]",
	Short: "Profile a table and summarize basic statistics",
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

		profile, err := analytics.ProfileTable(db, args[0])
		if err != nil {
			return err
		}
		displayProfile(profile)
		return nil
	},
}

func displayProfile(profile *analytics.Profile) {
	if profile == nil {
		return
	}
	fmt.Println("DATA PROFILE")
	fmt.Printf("Rows: %d\n", profile.RowCount)
	fmt.Printf("Columns: %d\n\n", profile.ColumnCount)
	fmt.Println("Column         Type            Missing   Unique   Min        Max        Mean   Median")
	fmt.Println("-------------------------------------------------------------------------------------------")
	for _, col := range profile.Columns {
		fmt.Printf("%-14s %-15s %-8.2f %-7d %-10s %-10s %-6s %-6s\n",
			col.Name,
			col.Type,
			col.Missing,
			col.Unique,
			truncate(col.Min, 10),
			truncate(col.Max, 10),
			truncate(col.Mean, 6),
			truncate(col.Median, 6),
		)
	}
}

func truncate(value string, width int) string {
	if value == "" {
		return ""
	}
	if len(value) <= width {
		return value
	}
	return value[:width-1] + "…"
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	_ = os.Stdout
	_ = strings.TrimSpace
}
