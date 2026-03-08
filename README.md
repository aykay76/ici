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

- Go 1.25 or later
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

### Environment Variables & Secrets

Store and manage secrets locally for use in workflows:

```bash
# Store a secret
ici secrets set MY_API_KEY sk-api-key-12345

# List all stored secrets
ici secrets list

# Remove a secret
ici secrets remove MY_API_KEY

# Run workflow with secrets (uses ~/.ici/secrets.json by default)
ici run .github/workflows/test.yml

# Use custom secrets file location
ici run .github/workflows/test.yml --secrets ./my-secrets.json
```

Secrets are stored in plain text in `~/.ici/secrets.json` (read/write 0600 permissions). They are automatically injected as environment variables into your workflows. This allows you to:

- Define workflow-level env vars: `env:` in the workflow YAML
- Define job-level env vars: `env:` in the job
- Define step-level env vars: `env:` in the step
- Inject stored secrets as environment variables at runtime

**Note:** Secrets are local to your machine and not synced with GitHub. They are useful for development/testing workflows locally.

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
│   │   ├── validate.go   # Validate command
│   │   └── secrets.go    # Secrets management command
│   ├── parser/           # Workflow parsing
│   │   └── workflow.go   # YAML parser & types
│   ├── runner/           # Workflow execution
│   │   └── executor.go   # Job & step execution
│   ├── container/        # Container management
│   │   └── podman.go     # Podman integration
│   └── secrets/          # Secret storage
│       └── store.go      # Local file-based secret store
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
- [x] Basic step execution (run steps, output capture, timeouts, exit codes)
- [x] Ubuntu container support
- [x] Environment variables (workflow, job, step levels)
- [x] Local secret management & injection
- [x] Working directory support with container mounts

### Phase 2: Action Support
- [ ] Action resolution
- [ ] actions/checkout implementation
- [ ] Action caching
- [ ] Working directory in steps

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
