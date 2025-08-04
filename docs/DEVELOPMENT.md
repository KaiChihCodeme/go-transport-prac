# Development Guide

## Development Environment Setup

This guide covers setting up your development environment for the Go Transport Practice project.

## Prerequisites

- **Go**: Version 1.21 or higher
- **Docker**: For containerized development (optional)
- **Make**: For build automation (optional but recommended)
- **Git**: For version control

### Installing Prerequisites

#### Go Installation
```bash
# On macOS with Homebrew
brew install go

# On Ubuntu/Debian
sudo apt-get install golang-go

# On Windows
# Download from https://golang.org/dl/
```

#### Docker Installation
```bash
# On macOS with Homebrew
brew install docker docker-compose

# On Ubuntu/Debian
sudo apt-get install docker.io docker-compose

# On Windows
# Download Docker Desktop from https://www.docker.com/products/docker-desktop
```

## Local Development Setup

### 1. Clone and Initialize

```bash
git clone <repository-url>
cd go-transport-prac
make deps
make tidy
```

### 2. Verify Setup

```bash
# Check Go version
go version

# Run tests to verify everything works
make test

# Build the project
make build
```

### 3. Development Workflow

```bash
# Format code before committing
make fmt

# Run static analysis
make vet

# Run linter (requires golangci-lint)
make lint

# Run all checks
make dev
```

## Docker Development

### 1. Start Development Environment

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps
```

### 2. Access Development Container

```bash
# Enter the Go development container
docker-compose exec go-dev /bin/bash

# Inside container, run commands
make test
make build
```

### 3. Available Services

| Service | Port | Purpose | Credentials |
|---------|------|---------|-------------|
| Go Dev | 8080-8082, 9090, 6060 | Development environment | - |
| PostgreSQL | 5432 | Database examples | user: transport_user, pass: transport_pass, db: transport_db |
| Redis | 6379 | Caching examples | - |
| MinIO | 9000, 9001 | Object storage examples | user: minioadmin, pass: minioadmin |

### 4. Managing Services

```bash
# Stop all services
docker-compose down

# Restart specific service
docker-compose restart go-dev

# View logs
docker-compose logs go-dev

# Rebuild containers
docker-compose up --build
```

## IDE Setup

### Visual Studio Code

Recommended extensions:
- **Go** (official Go extension)
- **Docker** (for container management)
- **REST Client** (for API testing)
- **Protocol Buffers** (for .proto files)

Settings for `.vscode/settings.json`:
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports",
    "go.testFlags": ["-v"],
    "go.coverOnSave": true
}
```

### GoLand/IntelliJ IDEA

1. Install Go plugin
2. Configure Go SDK to your installed version
3. Enable Go modules support
4. Set up run configurations for different components

## Development Tools

### Required Tools

```bash
# Install development tools
make install-tools

# Install golangci-lint for linting
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.55.2
```

### Protocol Buffers Tools

```bash
# Install protoc compiler
# On macOS
brew install protobuf

# On Ubuntu/Debian
sudo apt-get install protobuf-compiler

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./pkg/sdl/jsonschema/...

# Run benchmarks
make bench
```

### Test Structure

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test component interactions
- **Benchmark tests**: Performance testing
- **End-to-end tests**: Full workflow testing

## Debugging

### Local Debugging

```bash
# Run with race detection
go run -race ./cmd/main.go

# Debug with dlv
dlv debug ./cmd/main.go
```

### Container Debugging

```bash
# Access container with debugging tools
docker-compose exec go-dev /bin/bash

# Install debugging tools in container
go install github.com/go-delve/delve/cmd/dlv@latest
```

## Performance Profiling

### CPU Profiling

```bash
# Generate CPU profile
go test -cpuprofile=cpu.out -bench=.

# View profile
go tool pprof cpu.out
```

### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.out -bench=.

# View profile
go tool pprof mem.out
```

## Common Development Tasks

### Adding New SDL Technology

1. Create package under `pkg/sdl/newtechnology/`
2. Implement common interfaces
3. Add tests and benchmarks
4. Update documentation
5. Add example usage

### Adding New Web Protocol

1. Create package under `pkg/webprotocol/newprotocol/`
2. Implement service interfaces
3. Add client and server examples
4. Write comprehensive tests
5. Document usage patterns

### Creating Examples

1. Add to `examples/` directory
2. Include README with setup instructions
3. Provide sample data in `testdata/`
4. Document in main documentation

## Troubleshooting

### Common Issues

1. **Module not found**: Run `go mod tidy`
2. **Port conflicts**: Check `docker-compose ps` and adjust ports
3. **Permission issues**: Check Docker daemon permissions
4. **Build failures**: Verify Go version and dependencies

### Getting Help

1. Check project documentation in `docs/`
2. Review examples in `examples/`
3. Check specification documents in `llms/specs/`
4. Create issue with detailed error information

## Code Style Guidelines

### Go Conventions

- Follow official Go style guide
- Use `gofmt` for formatting
- Write meaningful variable names
- Include package documentation
- Write tests for all public functions

### Commit Messages

```
type(scope): brief description

Longer description if needed

- Bullet points for details
- Reference issues: #123
```

Types: feat, fix, docs, style, refactor, test, chore

### Documentation

- Keep README files updated
- Document all public APIs
- Include usage examples
- Maintain architecture documentation