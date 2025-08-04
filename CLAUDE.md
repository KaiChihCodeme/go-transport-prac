# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go transport practice project focused on learning and implementing common transportation technologies and protocols. The project serves as a learning repository with examples and implementations across two main areas:

**Schema Definition Languages (SDL):**
- JSON Schema
- Protobuf  
- Parquet
- Avro

**Web Protocols:**
- RESTful API
- gRPC
- WebSocket
- WebRTC
- GraphQL

## Project Structure

The project is organized with planned directories:
- `../sdl/` - Schema Definition Language examples and implementations
- `../web-protocol/` - Web protocol examples and implementations

Currently, the repository is in initial setup phase with only basic project files.

## Development Commands

Since this is a new Go project without established build configuration:

- **Initialize module**: `go mod init go-transport-prac`
- **Run**: `go run .` or `go run main.go`
- **Build**: `go build`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Tidy dependencies**: `go mod tidy`

## Specification Documents

The project includes comprehensive specification documents in `llms/specs/`:

- **requirements.md**: Detailed functional and technical requirements for all technologies
- **design.md**: Architecture design with library selections and implementation strategies  
- **tasks.md**: Development roadmap with phases, tasks, and milestones

These documents provide the complete blueprint for implementing the project using mainstream libraries and best practices.

## Architecture Notes

This is a learning-focused repository designed to provide practical examples of transportation technologies in Go. Each technology area (SDL and web protocols) will have separate example implementations demonstrating different approaches and use cases, following the modular architecture defined in the specification documents.
