# ici - Local GitHub Actions Runner

A Podman-based local GitHub Actions runner with AI-powered pre-flight analysis to catch issues before they reach remote CI.

## Overview

`ici` (pronounced "icky") allows you to run GitHub Actions workflows locally in Podman containers, providing fast feedback during development without waiting for remote CI pipelines.

## Features

- ✅ Parse GitHub Actions workflow files
- ✅ Execute workflows in isolated Podman containers
- ✅ Support for ubuntu-latest runners
- 🚧 Action resolution and caching
- 🚧 Private repository support
- 🚧 AI-powered pre-flight analysis

## Installation

### Prerequisites

- Go 1.23 or later
- Podman installed and configured

### Build from Source

```bash
# Clone the repository
git clone https://github.com/aykay76/ici.git
cd ici

# Download dependencies
go mod download

# Build the binary
go build -o ici ./cmd/ici

# Optionally install to your PATH
sudo mv ici /usr/local/bin/
```

## Usage

### Run a Workflow

Execute a GitHub Actions workflow file locally:

```bash
# Run all jobs in a workflow
ici run .github/workflows/test.yml

# Run a specific job
ici run .github/workflows/build.yml --job build

# Dry run (parse without executing)
ici run workflow.yml --dry-run

# Verbose output
ici run workflow.yml -v
```

### Parse a Workflow

Parse and display workflow structure:

```bash
# Parse and display as YAML
ici parse .github/workflows/test.yml

# Parse and display as JSON
ici parse workflow.yml --format json
```

### Validate a Workflow

Check workflow syntax and structure:

```bash
# Basic validation
ici validate .github/workflows/test.yml

# Strict validation
ici validate workflow.yml --strict
```

## Project Structure

```
ici/
├── cmd/
│   └── ici/              # CLI entry point
│       └── main.go
├── internal/
│   ├── cmd/              # CLI commands
│   │   ├── root.go       # Root command
│   │   ├── run.go        # Run command
│   │   ├── parse.go      # Parse command
│   │   └── validate.go   # Validate command
│   ├── parser/           # Workflow parsing
│   │   └── workflow.go   # YAML parser & types
│   ├── runner/           # Workflow execution
│   │   └── executor.go   # Job & step execution
│   └── container/        # Container management
│       └── podman.go     # Podman integration
├── go.mod
└── README.md
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Development

```bash
go build -o ici ./cmd/ici
./ici --help
```

## Roadmap

### Phase 1: Basic Runner (Current)
- [x] CLI scaffolding
- [x] Workflow parser
- [ ] Basic step execution
- [ ] Ubuntu container support
- [ ] Simple run commands

### Phase 2: Action Support
- [ ] Action resolution
- [ ] actions/checkout implementation
- [ ] Action caching
- [ ] Environment variables
- [ ] Working directory support

### Phase 3: Advanced Features
- [ ] Private repository support
- [ ] Multi-job dependencies
- [ ] Matrix builds
- [ ] Artifacts & caching
- [ ] Service containers

### Phase 4: AI Integration
- [ ] Pre-flight analysis
- [ ] Breaking change detection
- [ ] Security scanning
- [ ] Performance optimization

## Contributing

Contributions are welcome! Please follow these guidelines:

- Write tests for new features
- Follow Go best practices
- Keep dependencies minimal
- Focus on one task at a time

## License

MIT License - see LICENSE file for details

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing
- [Podman](https://podman.io/) - Container runtime
