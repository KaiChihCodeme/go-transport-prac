# Apache Avro Implementation in Go

This package provides a comprehensive implementation of Apache Avro serialization in Go, demonstrating various features including schema definition, JSON/binary encoding, schema evolution, schema registry concepts, and performance comparisons.

## Overview

Apache Avro is a data serialization system that provides:
- Rich data structures
- Schema evolution
- Compact binary data format
- Integration with dynamic languages

This implementation showcases all major Avro features using the `github.com/hamba/avro/v2` library.

## Features

- ✅ **Schema Definition**: Complete Avro schema definitions for User, Product, and Order entities
- ✅ **JSON Encoding**: Avro JSON serialization with schema validation
- ✅ **Binary Encoding**: Compact binary serialization for efficient storage/transmission
- ✅ **Schema Evolution**: Forward/backward compatibility with version management
- ✅ **Schema Registry**: Simulated schema registry with compatibility checking
- ✅ **Performance Comparison**: Benchmarks against standard JSON serialization
- ✅ **File Operations**: Reading/writing Avro files
- ✅ **Data Validation**: Schema-based validation
- ✅ **Comprehensive Tests**: Full test coverage with examples

## Quick Start

### Installation

```bash
go mod tidy
```

### Basic Usage

```go
package main

import (
    "fmt"
    "go-transport-prac/pkg/sdl/avro"
)

func main() {
    // Create manager
    manager, err := avro.NewManager("data/avro")
    if err != nil {
        panic(err)
    }

    // Create sample user
    user := avro.User{
        ID:     1,
        Email:  "user@example.com",
        Name:   "John Doe",
        Status: avro.UserStatusActive,
    }

    // Serialize to JSON
    jsonData, err := manager.SerializeUserJSON(user)
    if err != nil {
        panic(err)
    }

    // Serialize to binary
    binaryData, err := manager.SerializeUserBinary(user)
    if err != nil {
        panic(err)
    }

    fmt.Printf("JSON size: %d bytes\n", len(jsonData))
    fmt.Printf("Binary size: %d bytes\n", len(binaryData))
}
```

### Running Examples

Run all examples with:

```bash
go run cmd/avro_examples/main.go
```

This will demonstrate:
- JSON and binary encoding/decoding
- File operations
- Schema introspection
- Data validation
- Schema evolution
- Schema registry concepts
- Performance comparisons

## Schema Definitions

### User Schema (user.avsc)

```json
{
  "type": "record",
  "name": "User",
  "namespace": "com.example.avro",
  "fields": [
    {"name": "id", "type": "long"},
    {"name": "email", "type": "string"},
    {"name": "name", "type": "string"},
    {"name": "status", "type": {"type": "enum", "name": "UserStatus", "symbols": ["ACTIVE", "INACTIVE", "SUSPENDED", "DELETED"]}},
    {"name": "profile", "type": ["null", "Profile"]},
    {"name": "createdAt", "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name": "updatedAt", "type": {"type": "long", "logicalType": "timestamp-millis"}}
  ]
}
```

### Product Schema (product.avsc)

Complete schema with nested Price and Inventory records, arrays, maps, and optional fields.

### Order Schema (order.avsc)

Complex schema demonstrating nested records, arrays of records, and multiple optional fields.

## API Reference

### Manager

```go
type Manager struct {
    // Core Avro serialization manager
}

// Create new manager
func NewManager(baseDir string) (*Manager, error)

// JSON Serialization
func (m *Manager) SerializeUserJSON(user User) ([]byte, error)
func (m *Manager) DeserializeUserJSON(data []byte) (User, error)

// Binary Serialization  
func (m *Manager) SerializeUserBinary(user User) ([]byte, error)
func (m *Manager) DeserializeUserBinary(data []byte) (User, error)

// File Operations
func (m *Manager) WriteUsersToFile(filename string, users []User) error
func (m *Manager) ReadUsersFromFile(filename string) ([]User, error)

// Schema Access
func (m *Manager) GetUserSchema() avro.Schema
func (m *Manager) GetProductSchema() avro.Schema
func (m *Manager) GetOrderSchema() avro.Schema

// Sample Data
func (m *Manager) CreateSampleUsers(count int) []User
func (m *Manager) CreateSampleProducts(count int) []Product
```

### Schema Evolution

```go
type EvolutionManager struct {
    // Manages schema evolution scenarios
}

func NewEvolutionManager(baseDir string) (*EvolutionManager, error)
func (em *EvolutionManager) DemonstrateSchemaEvolution() error
func (em *EvolutionManager) GetSchemaVersions() map[string]string
```

### Schema Registry

```go
type SchemaRegistry struct {
    // Simulated schema registry
}

func NewSchemaRegistry() *SchemaRegistry
func (sr *SchemaRegistry) RegisterSchema(subject string, schemaJSON string) (int, error)
func (sr *SchemaRegistry) GetLatestSchema(subject string) (SchemaMetadata, error)
func (sr *SchemaRegistry) SetCompatibilityLevel(subject string, level CompatibilityLevel) error
```

## Schema Evolution

This implementation demonstrates three schema versions:

### v1 (Original)
- Basic user fields: id, email, name, status, profile
- Timestamps: createdAt, updatedAt

### v2 (Enhanced)
- Added optional fields: dateOfBirth, preferredLanguage, lastLoginAt
- Maintains backward compatibility with defaults

### v3 (Extended)
- Added enum value: ARCHIVED status
- Added nested fields: coordinates in address
- Added derived field: fullName in profile

### Evolution Rules

**Forward Compatibility:**
- Add new fields with default values
- Don't remove or rename existing fields
- Don't change field types
- Add new enum symbols at the end

**Backward Compatibility:**
- Make new fields optional (union with null)
- Provide sensible default values
- Don't remove enum symbols
- Consider aliases for field renames

## Performance Comparison

Benchmark results (1000 items):

| Format | Serialization | Deserialization | Size (bytes) | Memory (KB) | Items/sec |
|--------|---------------|-----------------|--------------|-------------|-----------|
| Avro JSON | ~8ms | ~8ms | 197 | 11,924 | 124,945 |
| Avro Binary | ~7ms | ~7ms | 197 | 13,994 | 149,341 |
| Standard JSON | ~3ms | ~3ms | 477 | 2,255 | 374,012 |

### Key Findings

- **Standard JSON**: Fastest serialization/deserialization, lowest memory usage
- **Avro JSON/Binary**: ~58% smaller serialized size, schema validation
- **Avro Binary**: Most compact, schema evolution support
- **Trade-offs**: Avro provides schema validation and evolution at performance cost

## Testing

Run tests with:

```bash
# Run all Avro tests
go test ./pkg/sdl/avro/... -v

# Run specific test
go test ./pkg/sdl/avro/... -v -run TestAvroManagerCreation

# Run benchmarks
go test ./pkg/sdl/avro/... -bench=. -benchmem
```

## Use Cases

### When to Use Avro

✅ **Good for:**
- Schema evolution requirements
- Cross-language data exchange
- Compact binary serialization
- Schema validation
- Big data processing (Kafka, Hadoop)

❌ **Consider alternatives for:**
- Simple data structures
- Human-readable formats required
- Maximum performance critical
- No schema evolution needed

### Comparison with Other Formats

| Feature | Avro | JSON | Protobuf | Parquet |
|---------|------|------|----------|---------|
| Schema Evolution | ✅ Excellent | ❌ None | ⚠️ Limited | ✅ Good |
| Binary Size | ✅ Compact | ❌ Large | ✅ Very Compact | ✅ Columnar |
| Human Readable | ⚠️ JSON variant | ✅ Yes | ❌ No | ❌ No |
| Language Support | ✅ Wide | ✅ Universal | ✅ Wide | ⚠️ Limited |
| Validation | ✅ Schema-based | ❌ None | ✅ Schema-based | ✅ Schema-based |

## File Structure

```
pkg/sdl/avro/
├── README.md              # This file
├── schemas/               # Avro schema definitions
│   ├── user.avsc         # User entity schema (v1)
│   ├── user_v2.avsc      # User schema v2 (evolution)
│   ├── user_v3.avsc      # User schema v3 (evolution)
│   ├── product.avsc      # Product entity schema
│   └── order.avsc        # Order entity schema
├── models.go              # Go struct definitions
├── manager.go             # Core serialization manager
├── converters.go          # Avro map conversion utilities
├── examples.go            # Usage examples and demonstrations
├── evolution.go           # Schema evolution examples
├── registry.go            # Schema registry simulation
├── benchmark.go           # Performance comparison benchmarks
└── manager_test.go        # Comprehensive tests
```

## Best Practices

### Schema Design

1. **Use logical types** for timestamps, decimals, UUIDs
2. **Make fields optional** when possible (union with null)
3. **Provide default values** for forward compatibility
4. **Use enums** for controlled vocabularies
5. **Document schemas** with meaningful names and documentation

### Evolution Strategy

1. **Plan for evolution** from day one
2. **Use semantic versioning** for schemas
3. **Test compatibility** before deployment
4. **Maintain schema registry** for centralized management
5. **Document breaking changes** clearly

### Performance Optimization

1. **Reuse schema objects** - parse once, use many times
2. **Batch operations** for file I/O
3. **Use binary format** for storage/transmission
4. **Profile your use case** - performance varies by data structure
5. **Consider compression** for additional space savings

## Integration Examples

### With Kafka

```go
// Producer
producer := kafka.NewProducer(config)
avroData, _ := manager.SerializeUserBinary(user)
producer.Produce(&kafka.Message{
    TopicPartition: kafka.TopicPartition{Topic: &topic},
    Value: avroData,
})

// Consumer  
for msg := range consumer.Events() {
    user, _ := manager.DeserializeUserBinary(msg.Value)
    // Process user...
}
```

### With File Storage

```go
// Write batch of users
users := manager.CreateSampleUsers(1000)
err := manager.WriteUsersToFile("users_batch.avro", users)

// Read batch back
readUsers, err := manager.ReadUsersFromFile("users_batch.avro")
```

### With HTTP API

```go
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    
    // Validate with Avro schema
    avroData, err := manager.SerializeUserJSON(user)
    if err != nil {
        http.Error(w, "Invalid user data", 400)
        return
    }
    
    // Store or process...
}
```

## References

- [Apache Avro Specification](https://avro.apache.org/docs/current/spec.html)
- [hamba/avro Go Library](https://github.com/hamba/avro)
- [Schema Evolution Best Practices](https://docs.confluent.io/platform/current/schema-registry/avro.html)
- [Avro vs Other Formats](https://avro.apache.org/docs/current/index.html)

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Run performance benchmarks
5. Validate schema compatibility

---

This implementation provides a complete foundation for working with Apache Avro in Go, from basic serialization to advanced schema evolution and registry concepts.