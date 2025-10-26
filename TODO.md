# ici - Development Backlog

This document tracks the development roadmap and tasks for the `ici` (ShiftLeft CI) project.

## 🎯 Current Status
- ✅ CLI scaffolding complete
- ✅ Basic workflow parser (YAML)
- ✅ Command structure (run, parse, validate)
- ✅ Stub execution with verbose output

---

## 📋 Phase 1: Basic Runner (In Progress)

### High Priority

- [ ] **Podman Integration**
  - [x] Implement `CreateContainer()` in `internal/container/podman.go` (now returns container ID)
  - [ ] Pull container images (ubuntu:22.04, etc.)
  - [ ] Start containers with proper configuration
  - [x] Implement `RunCommand()` to execute shell commands in containers
  - [x] Implement `RemoveContainer()` for cleanup
  - [ ] Handle container lifecycle (start, stop, remove)

- [ ] **Basic Step Execution**
  - [ ] Execute `run:` steps in containers
  - [ ] Capture stdout/stderr from container commands
  - [ ] Handle step failures and exit codes
  - [ ] Display step output in real-time
  - [ ] Implement step timeout handling

Note: Unit tests for the container manager were added (tests stub CLI behavior and validate create/exec/remove). Integration tests remain as a follow-up.

- [ ] **Environment Variables**
  - [ ] Parse workflow-level `env:`
  - [ ] Parse job-level `env:`
  - [ ] Parse step-level `env:`
  - [ ] Pass environment variables to containers
  - [ ] Support for secret handling (placeholder for now)

- [ ] **Working Directory Support**
  - [ ] Mount workspace directory into containers
  - [ ] Handle `working-directory` in steps
  - [ ] Ensure proper path mapping between host and container

### Medium Priority

- [ ] **Error Handling & Logging**
  - [ ] Structured logging (consider zerolog or similar)
  - [ ] Better error messages with context
  - [ ] Exit codes that match GitHub Actions behavior
  - [ ] Log file output option

- [ ] **Testing**
  - [ ] Unit tests for parser package
  - [ ] Unit tests for runner package
  - [ ] Integration tests with actual Podman
  - [ ] Test fixtures (sample workflows)

- [ ] **actions/checkout Implementation**
  - [ ] Clone current repository into workspace
  - [ ] Handle different refs (branches, tags, commits)
  - [ ] Sparse checkout support
  - [ ] Submodule handling

---

## 📋 Phase 2: Action Support

### High Priority

- [ ] **Action Resolution**
  - [ ] Parse `uses:` syntax (owner/repo@ref)
  - [ ] Download actions from GitHub
  - [ ] Cache downloaded actions locally
  - [ ] Handle action versioning (tags, SHAs, branches)

- [ ] **Composite Actions**
  - [ ] Parse `action.yml` files
  - [ ] Execute composite action steps
  - [ ] Handle action inputs/outputs
  - [ ] Support for nested actions

- [ ] **Docker Actions**
  - [ ] Build Docker-based actions
  - [ ] Execute in separate containers
  - [ ] Handle Dockerfile actions

### Medium Priority

- [ ] **JavaScript Actions**
  - [ ] Set up Node.js environment
  - [ ] Execute `node` actions
  - [ ] Handle action dependencies (npm install)

- [ ] **Action Outputs**
  - [ ] Capture step outputs
  - [ ] Make outputs available to subsequent steps
  - [ ] Support `${{ steps.id.outputs.name }}` syntax

- [ ] **Common Actions**
  - [ ] actions/setup-node
  - [ ] actions/setup-python
  - [ ] actions/setup-go
  - [ ] actions/cache
  - [ ] actions/upload-artifact
  - [ ] actions/download-artifact

---

## 📋 Phase 3: Advanced Workflow Features

### High Priority

- [ ] **Job Dependencies**
  - [ ] Parse and respect `needs:` relationships
  - [ ] Build job dependency graph
  - [ ] Execute jobs in correct order
  - [ ] Handle job failures in dependency chain

- [ ] **Matrix Builds**
  - [ ] Parse `strategy.matrix`
  - [ ] Generate job instances from matrix
  - [ ] Execute matrix jobs in parallel (configurable)
  - [ ] Display matrix results clearly

- [ ] **Conditionals**
  - [ ] Implement expression evaluation for `if:`
  - [ ] Support GitHub Actions context variables
  - [ ] Job-level conditionals
  - [ ] Step-level conditionals

### Medium Priority

- [ ] **Artifacts & Caching**
  - [ ] Local artifact storage
  - [ ] Upload/download artifacts between jobs
  - [ ] Cache implementation (paths, keys)
  - [ ] Cache restore/save

- [ ] **Service Containers**
  - [ ] Parse `services:` in jobs
  - [ ] Start service containers
  - [ ] Network configuration between containers
  - [ ] Health checks

- [ ] **Secrets & Variables**
  - [ ] Read from `.env` file
  - [ ] Command-line secret passing
  - [ ] Secure secret handling in containers
  - [ ] GitHub Variables support

---

## 📋 Phase 4: Authentication & Private Repos

### High Priority

- [ ] **Git Credentials**
  - [ ] Detect SSH keys from local agent
  - [ ] Use PAT tokens from git config
  - [ ] Pass credentials securely to containers
  - [ ] Support for `.netrc` files

- [ ] **Private Repositories**
  - [ ] Clone private repos in actions/checkout
  - [ ] Access private actions
  - [ ] Handle GitHub API rate limits

### Medium Priority

- [ ] **Docker Registry Authentication**
  - [ ] Support for private Docker registries
  - [ ] GitHub Container Registry (ghcr.io)
  - [ ] Docker Hub authentication

---

## 📋 Phase 5: Multi-Platform Support

### Future Enhancements

- [ ] **Runner OS Support**
  - [ ] macOS support (via Docker)
  - [ ] Windows support (via Docker)
  - [ ] Container image mapping for different OS

- [ ] **Self-Hosted Runner Compatibility**
  - [ ] Match self-hosted runner behavior
  - [ ] Custom labels support

---

## 📋 Phase 6: AI-Powered Features

### Future Enhancements

- [ ] **Pre-flight Analysis**
  - [ ] Detect outdated action versions
  - [ ] Suggest action updates
  - [ ] Security vulnerability scanning
  - [ ] Breaking change detection

- [ ] **Performance Optimization**
  - [ ] Identify slow steps
  - [ ] Suggest caching strategies
  - [ ] Cost estimation (if running on cloud)

- [ ] **Best Practices**
  - [ ] Workflow linting
  - [ ] Suggest improvements
  - [ ] Detect anti-patterns

---

## 📋 Developer Experience

### Ongoing

- [ ] **Documentation**
  - [ ] API documentation (godoc)
  - [ ] User guide with examples
  - [ ] Troubleshooting guide
  - [ ] Contributing guidelines

- [ ] **CLI Improvements**
  - [ ] Shell completion (bash, zsh, fish)
  - [ ] Better progress indicators
  - [ ] Colored output
  - [ ] Interactive mode

- [ ] **Configuration**
  - [ ] Config file support (.ici.yml)
  - [ ] Per-repository configuration
  - [ ] Global user configuration

---

## 🐛 Known Issues

- None yet! This is a fresh project.

---

## 💡 Ideas for Future Consideration

- [ ] Web UI for viewing run results
- [ ] GitHub App integration
- [ ] Slack/Discord notifications
- [ ] Workflow visualization (graph of jobs/steps)
- [ ] Diff mode (compare local vs remote CI results)
- [ ] Record and replay mode
- [ ] Plugin system for custom actions
- [ ] Performance profiling mode
- [ ] Integration with CI/CD platforms beyond GitHub Actions

---

## 📝 Notes

- Keep dependencies minimal (per `.github/copilotinstructions.md`)
- Focus on one task at a time
- Write tests as features are implemented
- Maintain compatibility with GitHub Actions syntax
