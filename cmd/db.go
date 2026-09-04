package main

import (
	"fmt"
	"os"
	"strings"

	"datanex/internal/config"

	"github.com/spf13/cobra"
)

var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured database connections",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if len(cfg.Connections) == 0 {
			fmt.Println("No database connections configured.")
			return nil
		}
		fmt.Println("Configured connections:")
		for _, conn := range cfg.Connections {
			marker := " "
			if cfg.DefaultDB == conn.Name {
				marker = "*"
			}
			fmt.Printf("%s %s (%s)\n", marker, conn.Name, conn.Type)
		}
		return nil
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [name] [type] [dsn]",
	Short: "Add or update a database connection",
	Args:  cobra.RangeArgs(1, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		name := args[0]
		connType := "sqlite"
		if len(args) >= 2 {
			connType = args[1]
		}
		dsn := "file:default.db"
		if len(args) >= 3 {
			dsn = args[2]
		}

		conn := config.Connection{Name: name, Type: connType, DSN: dsn}
		updated := false
		for i := range cfg.Connections {
			if cfg.Connections[i].Name == name {
				cfg.Connections[i] = conn
				updated = true
				break
			}
		}
		if !updated {
			cfg.Connections = append(cfg.Connections, conn)
		}
		cfg.DefaultDB = name
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("save connection: %w", err)
		}
		fmt.Printf("Connection %q saved successfully (%s).\n", name, connType)
		return nil
	},
}

var testConnectionCmd = &cobra.Command{
	Use:   "test [name]",
	Short: "Test a configured connection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		conn, err := findConnection(cfg, args[0])
		if err != nil {
			return err
		}
		if strings.EqualFold(conn.Type, "sqlite") {
			db, err := openDBConnection(conn)
			if err != nil {
				return fmt.Errorf("test connection: %w", err)
			}
			defer db.Close()
			fmt.Printf("Connection %q is healthy.\n", conn.Name)
			return nil
		}
		if strings.EqualFold(conn.Type, "postgres") || strings.EqualFold(conn.Type, "postgresql") {
			return fmt.Errorf("postgresql support is implemented in a later phase")
		}
		return fmt.Errorf("unsupported database type %q", conn.Type)
	},
}

var removeConnectionCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a stored connection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		filtered := make([]config.Connection, 0, len(cfg.Connections))
		removed := false
		for _, conn := range cfg.Connections {
			if conn.Name == args[0] {
				removed = true
				continue
			}
			filtered = append(filtered, conn)
		}
		if !removed {
			return fmt.Errorf("connection %q not found", args[0])
		}
		cfg.Connections = filtered
		if cfg.DefaultDB == args[0] {
			cfg.DefaultDB = ""
		}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("save updated config: %w", err)
		}
		fmt.Printf("Connection %q removed.\n", args[0])
		return nil
	},
}

func init() {
	dbCmd.AddCommand(listConnectionsCmd, connectCmd, testConnectionCmd, removeConnectionCmd)
	_ = os.Stdout
}
