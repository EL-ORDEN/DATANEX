# DataNex

DataNex is a modern Go-based command-line data exploration and management tool built for practical use in local data engineering workflows. It is designed to connect to databases, execute SQL, import and export data, query APIs, and profile datasets from a professional CLI.

## Overview

The project demonstrates a real-world architecture for a developer-facing data tool:

- database connectivity with SQLite and PostgreSQL
- SQL query execution and terminal exploration
- CSV and JSON import/export
- REST API testing with auth and headers
- basic data profiling and analytics
- YAML/JSON pipeline configuration
- modular, testable Go architecture

## Features

- Database connection management: connect, list, test, remove
- SQLite and PostgreSQL handling
- SQL query runner and interactive shell
- CSV/JSON import into a target table
- CSV/JSON export from tables
- REST API client with GET, POST, PUT, PATCH, DELETE
- Header and bearer-token support
- Basic table profiling and statistics
- Pipeline definitions via YAML/JSON
- Config persistence and secure redaction of connection strings

## Architecture

The code is organized into small, maintainable packages:

- `cmd/` – CLI entrypoints and commands
- `internal/config/` – app settings and secure connection storage
- `internal/database/` – SQLite/PostgreSQL connections and metadata helpers
- `internal/query/` – SQL execution engine
- `internal/importer/` – CSV/JSON import workflows
- `internal/exporter/` – CSV/JSON export workflows
- `internal/api/` – HTTP client and API interaction layer
- `internal/analytics/` – data profile and summary statistics
- `internal/pipeline/` – YAML/JSON pipeline definitions and validation
- `internal/ui/` – interactive terminal shell support

## Installation

Requirements:

- Go 1.25 or newer
- SQLite support through the Go driver
- PostgreSQL support through `pgx`

Install and build:

```bash
go mod tidy
go build ./...
```

## Usage

Show the CLI help:

```bash
go run ./cmd --help
```

Connect a database:

```bash
go run ./cmd db connect demo sqlite "file:demo.db"
```

List configured connections:

```bash
go run ./cmd db list
```

Run a SQL query:

```bash
go run ./cmd query "SELECT * FROM users LIMIT 10;"
```

Open the SQL shell:

```bash
go run ./cmd shell
```

Import CSV:

```bash
go run ./cmd import ./data/users.csv --table users --create=true
```

Export JSON:

```bash
go run ./cmd export users --format json --output ./out/users.json
```

Call an API:

```bash
go run ./cmd api get https://api.example.com/users --header "Accept:application/json"
```

Profile a table:

```bash
go run ./cmd analyze users
```

Run a pipeline:

```bash
go run ./cmd pipeline run ./examples/sample_pipeline.yaml
```

## Example commands

```bash
go run ./cmd db connect app sqlite "file:app.db"
go run ./cmd query "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER);"
go run ./cmd query "INSERT INTO users (name, age) VALUES ('Ada', 36), ('Linus', 54);"
go run ./cmd query "SELECT * FROM users;"
go run ./cmd api get https://jsonplaceholder.typicode.com/todos/1 --header "Accept:application/json"
```

## Configuration

DataNex stores settings in a per-user config directory:

- Windows: `%APPDATA%\DataNex\config.json`
- macOS/Linux: `~/.datanex/config.json`

Connections are persisted as JSON and passwords are redacted before display in logs or CLI output.

## Testing

Run the full automated test suite:

```bash
go test ./...
```

## Roadmap

Planned next milestones:

- richer PostgreSQL metadata explorer
- foreign key and index inspection
- better table rendering in the terminal
- more advanced analytics and distributions
- real pipeline execution steps for read/transform/validate/write
- logging and structured output improvements
- UX polish for interactive shell and command suggestions

## Contributing

Contributions are welcome. Please keep changes scoped, testable, and aligned with the modular architecture used in the project.

## License

MIT
