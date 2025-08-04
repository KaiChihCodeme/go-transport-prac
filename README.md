# go-transport-prac

This project is for practicing common and state-of-the-art transportation technologies in Go.
It provides examples and implementations to help understand these technologies better.

## Technologies Covered

### Schema Definition Languages (SDL)
Located in `pkg/sdl/`:
1. **JSON Schema** - API request/response validation
2. **Protocol Buffers** - High-performance serialization
3. **Parquet** - Columnar data storage
4. **Avro** - Schema evolution and streaming

### Web Protocols
Located in `pkg/webprotocol/`:
1. **RESTful API** - HTTP-based web services
2. **gRPC** - High-performance RPC
3. **WebSocket** - Real-time bidirectional communication
4. **WebRTC** - Peer-to-peer communication
5. **GraphQL** - Flexible query language

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose (optional)
- Make (optional, but recommended)

### Local Development

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd go-transport-prac
   make deps
   ```

2. **Run tests:**
   ```bash
   make test
   ```

3. **Format and lint:**
   ```bash
   make fmt
   make vet
   ```

### Docker Development

1. **Start development environment:**
   ```bash
   docker-compose up -d
   ```

2. **Access development container:**
   ```bash
   docker-compose exec go-dev /bin/bash
   ```

3. **Run commands inside container:**
   ```bash
   make test
   make build
   ```

### Available Services

When using Docker Compose, the following services are available:
- **Go Development**: Port 8080-8082, 9090, 6060
- **PostgreSQL**: Port 5432 (user: transport_user, db: transport_db)
- **Redis**: Port 6379
- **MinIO**: Port 9000 (console: 9001)

## Project Structure

```
go-transport-prac/
├── cmd/                    # CLI applications and demos
├── internal/               # Internal shared packages
├── pkg/                    # Public packages
│   ├── sdl/               # Schema Definition Languages
│   └── webprotocol/       # Web Protocols
├── examples/              # Standalone examples
├── docs/                  # Documentation
├── testdata/              # Test data and fixtures
└── llms/specs/            # Project specifications
```

## Development Commands

```bash
make help          # Show all available commands
make build         # Build the project
make test          # Run tests
make test-coverage # Run tests with coverage
make fmt           # Format code
make lint          # Run linter
make vet           # Run go vet
make clean         # Clean build artifacts
make deps          # Install dependencies
make tidy          # Tidy dependencies
make dev           # Development workflow (fmt, vet, test)
make docs          # Generate documentation
```

## Learning Path

1. **Start with SDL examples** to understand data serialization
2. **Explore Web Protocols** for communication patterns
3. **Check integration examples** in `examples/` directory
4. **Read documentation** in `docs/` for detailed explanations

## Contributing

This is a learning project. Feel free to explore, modify, and experiment with the code to better understand these transportation technologies.
