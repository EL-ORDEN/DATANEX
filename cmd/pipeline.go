package main

import (
	"fmt"
	"os"

	"datanex/internal/config"
	"datanex/internal/pipeline"

	"github.com/spf13/cobra"
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Run a data pipeline",
}

var pipelineRunCmd = &cobra.Command{
	Use:   "run [pipeline-file]",
	Short: "Run a YAML or JSON pipeline definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := pipeline.LoadFromFile(args[0])
		if err != nil {
			return err
		}
		appCfg, err := config.LoadConfig("")
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if appCfg.DefaultDB == "" {
			return fmt.Errorf("no default database configured")
		}
		conn, err := findConnection(appCfg, appCfg.DefaultDB)
		if err != nil {
			return err
		}
		db, err := openDBConnection(conn)
		if err != nil {
			return err
		}
		defer db.Close()
		if err := cfg.Run(db); err != nil {
			return err
		}
		fmt.Printf("Pipeline %q executed successfully.\n", cfg.Name)
		return nil
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineRunCmd)
	rootCmd.AddCommand(pipelineCmd)
	_ = os.Stdout
}
