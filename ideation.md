# Project: Local CI Runner ("ShiftLeft CI")

## Vision
Build a Podman-based local GitHub Actions runner with AI-powered pre-flight analysis to catch issues before they reach remote CI.

## Core Architecture

### Phase 1: Basic Ubuntu Runner
- **Input**: GitHub Actions YAML workflow files
- **Execution**: Podman containers (ubuntu-latest focus)
- **Output**: Local CI results with remote CI parity

### Key Technical Components:
1. YAML Parser → Execution Plan
2. Podman Container Orchestrator  
3. Action Resolver & Cache
4. Step Runner with State Management
5. Artifact/Cache Manager

## Current Focus: Workflow Parser & Basic Runner

### Immediate Tasks:
1. Parse GitHub Actions YAML format
2. Extract job definitions, steps, and container requirements
3. Map `runs-on` to appropriate Podman images
4. Execute simple `run` steps in containers
5. Handle basic environment variables and working directories

### Example Target Workflow Support:
```yaml
name: Test
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Run tests
        run: make test
```

Authentication Simplification Goals
Private Repo Access:
Auto-detect existing Git credentials

SSH key management from local agent

PAT token reuse from git config

Secure credential passing to containers

AI Integration Opportunities (Future)
Pre-flight Analysis:
Action version outdatedness detection

Breaking change prediction

Security vulnerability scanning

Performance bottleneck identification

Cost optimization suggestions

Development Principles
Technical Stack:
Language: Go (good Podman integration, performance)

Container: Podman (rootless, daemonless)

Parsing: Go YAML libraries + custom workflow semantics

Code Quality:
Test-driven development

Container-native design

CLI-first interface

Extensible architecture for future multi-platform support

Getting Started Commands
bash
# Project setup
git init local-ci-runner
cd local-ci-runner
go mod init github.com/yourusername/local-ci

# Basic structure
mkdir -p cmd/ internal/parser/ internal/runner/ internal/container/

# Initial implementation focus:
# 1. cmd/local-ci/main.go - CLI entrypoint
# 2. internal/parser/workflow.go - YAML parsing
# 3. internal/container/podman.go - Podman integration
# 4. internal/runner/executor.go - Step execution
Success Metrics
✅ Parse basic GitHub Actions YAML

✅ Run simple shell commands in Ubuntu container

✅ Handle actions/checkout equivalent

✅ Provide faster feedback than remote CI

✅ Work with private repositories seamlessly

Questions to Consider During Implementation:
How do we handle action version resolution?

What's the optimal container caching strategy?

How to securely pass credentials to containers?

What's the minimal feature set for MVP?

How to structure code for eventual multi-platform support?

Remember: Start simple, solve the 80% use case first, then expand.
