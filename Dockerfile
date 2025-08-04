# Development Dockerfile for Go Transport Practice

FROM golang:1.21-alpine AS development

# Install development tools
RUN apk add --no-cache \
    git \
    make \
    curl \
    bash \
    protobuf \
    protobuf-dev

# Set working directory
WORKDIR /app

# Install additional Go tools
RUN go install golang.org/x/tools/cmd/godoc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose ports for various services
EXPOSE 8080 8081 8082 9090 6060

# Default command for development
CMD ["sh", "-c", "make help && /bin/bash"]

# Production build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]