package ui

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"datanex/internal/query"
)

// SQLShell provides a minimal interactive shell for SQL execution.
type SQLShell struct {
	DB   *sql.DB
	Hist []string
}

// NewSQLShell creates a shell bound to a database connection.
func NewSQLShell(db *sql.DB) *SQLShell {
	return &SQLShell{DB: db}
}

// Run starts the interactive loop.
func (s *SQLShell) Run() error {
	fmt.Println("DataNex SQL shell")
	fmt.Println("Type 'exit' to quit.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("datanex> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println()
				return nil
			}
			return fmt.Errorf("read input: %w", err)
		}
		line := strings.TrimSpace(input)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			fmt.Println("Bye.")
			return nil
		}
		start := time.Now()
		result, err := query.NewEngine(s.DB).Execute(line)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		result.Duration = time.Since(start)
		fmt.Printf("Rows: %d | Time: %s\n", result.RowCount, result.Duration)
		if len(result.Columns) > 0 {
			fmt.Println(strings.Join(result.Columns, " | "))
			for _, row := range result.Rows {
				fmt.Println(strings.Join(row, " | "))
			}
		}
		s.Hist = append(s.Hist, line)
	}
}
