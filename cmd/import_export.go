package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"datanex/internal/config"
	"datanex/internal/exporter"
	"datanex/internal/importer"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file] --table [table]",
	Short: "Import CSV or JSON data into a table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, _ := cmd.Flags().GetString("table")
		create, _ := cmd.Flags().GetBool("create")
		replace, _ := cmd.Flags().GetBool("replace")
		if table == "" {
			return fmt.Errorf("table name is required; use --table")
		}
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

		file := args[0]
		switch strings.ToLower(filepath.Ext(file)) {
		case ".csv":
			count, err := importer.ImportCSV(db, file, table, create, replace)
			if err != nil {
				return err
			}
			fmt.Printf("Imported %d rows into %s from %s\n", count, table, file)
		case ".json":
			count, err := importer.ImportJSON(db, file, table, create, replace)
			if err != nil {
				return err
			}
			fmt.Printf("Imported %d rows into %s from %s\n", count, table, file)
		default:
			return fmt.Errorf("unsupported file type %q; use .csv or .json", filepath.Ext(file))
		}
		return nil
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [table] --format [csv|json]",
	Short: "Export a table to CSV or JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = fmt.Sprintf("%s.%s", args[0], strings.ToLower(format))
		}
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

		switch strings.ToLower(format) {
		case "csv":
			if err := exporter.ExportCSV(db, args[0], output); err != nil {
				return err
			}
		case "json":
			if err := exporter.ExportJSON(db, args[0], output); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported export format %q; use csv or json", format)
		}
		fmt.Printf("Exported %s to %s\n", args[0], output)
		return nil
	},
}

func init() {
	importCmd.Flags().String("table", "", "Target table name")
	importCmd.Flags().Bool("create", true, "Create the table if it does not exist")
	importCmd.Flags().Bool("replace", false, "Replace existing rows before insert")
	exportCmd.Flags().String("format", "csv", "Export format: csv or json")
	exportCmd.Flags().String("output", "", "Output path for the exported file")
	rootCmd.AddCommand(importCmd, exportCmd)
	_ = os.Stdout
}
