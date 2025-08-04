# Go Transport Practice Project Makefile

.PHONY: help build test fmt lint clean deps tidy vet run

# Default target
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the project
build: ## Build the project
	@echo "Building..."
	go build -v ./...

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.55.2" && exit 1)
	golangci-lint run

# Vet code
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Clean build artifacts
clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	go clean ./...
	rm -f coverage.out coverage.html

# Install dependencies
deps: ## Download and install dependencies
	@echo "Installing dependencies..."
	go mod download

# Tidy dependencies
tidy: ## Tidy go modules
	@echo "Tidying dependencies..."
	go mod tidy

# Run development server (when main.go exists)
run: ## Run the main application
	@echo "Running application..."
	@if [ -f cmd/main.go ]; then \
		go run cmd/main.go; \
	elif [ -f main.go ]; then \
		go run main.go; \
	else \
		echo "No main.go file found in root or cmd/ directory"; \
	fi

# Development workflow
dev: fmt vet test ## Run development workflow (format, vet, test)

# Full check (for CI/CD)
check: fmt vet lint test ## Run all checks (format, vet, lint, test)

# Generate documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@which godoc > /dev/null || go install golang.org/x/tools/cmd/godoc@latest
	@echo "Documentation server will be available at http://localhost:6060"
	@echo "Run 'godoc -http=:6060' to start the documentation server"

# Benchmark tests
bench: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Install tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/godoc@latest
	@echo "To install golangci-lint, run:"
	@echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.55.2"