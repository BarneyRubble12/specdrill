# SpecDrill

SpecDrill is a Go-based tool for automatically generating and running test cases for REST APIs based on OpenAPI specifications.

## Features

- Parse OpenAPI specifications (YAML/JSON) from local files or URLs
- Generate test cases from API endpoints
- Execute tests against target APIs
- CLI interface for easy usage
- Modular architecture for extensibility

## Installation

```bash
go install github.com/BarneyRubble12/specdrill/cmd/cli@latest
```

## Usage

```bash
# Using a local OpenAPI specification file
specdrill --spec ./openapi.yaml

# Using a remote OpenAPI specification URL
specdrill --spec https://api.example.com/openapi.json

# Overriding the base URL for testing
specdrill --spec ./openapi.yaml --base-url https://staging-api.example.com
```

## Project Structure

```
specdrill/
├── cmd/
│   └── cli/                 # CLI entry point
├── internal/
│   ├── core/                # Core logic
│   │   ├── parser/          # OpenAPI parsing
│   │   ├── generator/       # Test case generation
│   │   ├── executor/        # Test execution
│   │   └── model/           # Core domain models
│   ├── infrastructure/      # HTTP client, logging, utils
│   └── web/                 # Future web adapter
├── testdata/                # Example OpenAPI files
├── go.mod
└── README.md
```

## Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build: `go build ./cmd/cli`
4. Run tests: `go test ./...`

## License

MIT 