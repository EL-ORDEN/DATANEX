# DataNex

DataNex is a modern command-line data exploration and management tool built in Go. It combines database connectivity, SQL querying, import/export workflows, and API exploration into a single professional CLI.

## Features

- SQLite and PostgreSQL connection management
- SQL execution and interactive shell
- CSV and JSON import/export
- REST API exploration with auth and headers
- Data profiling and basic analytics
- YAML/JSON pipeline configuration
- Structured logging and configuration management

## Architecture

The project is organized around a modular Go architecture:

- `cmd/` – CLI commands and root command
- `internal/config/` – configuration and secure connection storage
- `internal/database/` – database abstractions and drivers
- `internal/query/` – execution engine and result formatting
- `internal/importer/` – CSV/JSON import workflows
- `internal/exporter/` – CSV/JSON/SQL export
- `internal/api/` – HTTP client wrappers and REST actions
- `internal/analytics/` – profiling and statistics
- `internal/ui/` – terminal UI components when needed

## Installation

```bash
go build ./...
```

## Usage

```bash
datanex --help
datanex db list
datanex query "SELECT * FROM users LIMIT 10"
datanex shell
```

## Phase 1 status

This branch focuses on the first development phase:

- project structure
- CLI foundation
- configuration management
- SQLite support
- basic query execution

## License

MIT
