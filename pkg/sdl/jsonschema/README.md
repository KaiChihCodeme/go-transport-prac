# JSON Schema Validation

This package provides JSON Schema validation capabilities using the `xeipuuv/gojsonschema` library.

## Features

- JSON Schema validation for JSON strings and Go objects
- HTTP middleware for request validation
- Comprehensive error reporting
- Schema management (add, remove, list)
- Detailed validation results

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "go-transport-prac/pkg/sdl/jsonschema"
    "go-transport-prac/internal/logger"
)

func main() {
    // Create logger
    logger, _ := logger.New(&logger.Config{Level: "info"})
    
    // Create validator
    validator := jsonschema.NewXeipuuvValidator(logger)

    // Define schema
    schemaJSON := `{
        "type": "object",
        "properties": {
            "name": {"type": "string", "minLength": 1},
            "age": {"type": "integer", "minimum": 0, "maximum": 150},
            "email": {"type": "string", "format": "email"}
        },
        "required": ["name", "age"]
    }`

    // Add schema
    err := validator.AddSchemaJSON("person", schemaJSON)
    if err != nil {
        log.Fatal(err)
    }

    // Validate JSON string
    jsonData := `{"name": "John Doe", "age": 30, "email": "john@example.com"}`
    err = validator.ValidateJSON("person", jsonData)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation successful!")
    }

    // Validate Go object
    data := map[string]interface{}{
        "name": "Jane Smith",
        "age":  25,
        "email": "jane@example.com",
    }
    err = validator.ValidateData("person", data)
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Object validation successful!")
    }
}
```

### HTTP Middleware Usage

```go
package main

import (
    "net/http"

    "go-transport-prac/pkg/sdl/jsonschema"
    "go-transport-prac/internal/logger"
)

func main() {
    // Create validator and middleware
    logger, _ := logger.New(&logger.Config{Level: "info"})
    validator := jsonschema.NewXeipuuvValidator(logger)
    middleware := jsonschema.NewSimpleHTTPMiddleware(validator, logger)

    // Add user creation schema
    userSchema := `{
        "type": "object",
        "properties": {
            "username": {"type": "string", "minLength": 3, "maxLength": 20},
            "password": {"type": "string", "minLength": 8},
            "email": {"type": "string", "format": "email"}
        },
        "required": ["username", "password", "email"]
    }`
    validator.AddSchemaJSON("create-user", userSchema)

    // Create routes with validation
    mux := http.NewServeMux()
    
    // Protected endpoint with validation
    createUserHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"status": "user created"}`))
    })
    mux.Handle("/users", middleware.ValidateRequest("create-user")(createUserHandler))

    // Standalone validation endpoint
    mux.Handle("/validate/user", middleware.ValidationHandler("create-user"))

    http.ListenAndServe(":8080", mux)
}
```

### Detailed Validation Results

```go
// Get detailed validation results
result, err := validator.ValidateWithDetails("person", invalidData)
if err != nil {
    log.Fatal(err)
}

if !result.Valid {
    fmt.Printf("Validation failed with %d errors:\n", len(result.Errors))
    for _, validationErr := range result.Errors {
        fmt.Printf("- Field '%s': %s (value: %v)\n", 
            validationErr.InstanceLocation, 
            validationErr.Message, 
            validationErr.Value)
    }
}
```

## Schema Examples

### User Profile Schema

```json
{
    "type": "object",
    "properties": {
        "id": {"type": "integer", "minimum": 1},
        "username": {
            "type": "string",
            "minLength": 3,
            "maxLength": 20,
            "pattern": "^[a-zA-Z0-9_]+$"
        },
        "email": {"type": "string", "format": "email"},
        "profile": {
            "type": "object",
            "properties": {
                "firstName": {"type": "string", "minLength": 1},
                "lastName": {"type": "string", "minLength": 1},
                "age": {"type": "integer", "minimum": 13, "maximum": 120},
                "bio": {"type": "string", "maxLength": 500}
            },
            "required": ["firstName", "lastName"]
        },
        "preferences": {
            "type": "object",
            "properties": {
                "theme": {"type": "string", "enum": ["light", "dark", "auto"]},
                "notifications": {"type": "boolean"},
                "language": {"type": "string", "default": "en"}
            }
        },
        "tags": {
            "type": "array",
            "items": {"type": "string"},
            "uniqueItems": true,
            "maxItems": 10
        }
    },
    "required": ["id", "username", "email", "profile"]
}
```

### Product Catalog Schema

```json
{
    "type": "object",
    "properties": {
        "sku": {"type": "string", "pattern": "^[A-Z0-9-]+$"},
        "name": {"type": "string", "minLength": 1, "maxLength": 100},
        "description": {"type": "string", "maxLength": 1000},
        "price": {
            "type": "object",
            "properties": {
                "amount": {"type": "number", "minimum": 0},
                "currency": {"type": "string", "enum": ["USD", "EUR", "GBP", "JPY"]}
            },
            "required": ["amount", "currency"]
        },
        "category": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "name": {"type": "string"},
                "path": {"type": "array", "items": {"type": "string"}}
            },
            "required": ["id", "name"]
        },
        "attributes": {
            "type": "object",
            "patternProperties": {
                "^[a-z_]+$": {
                    "oneOf": [
                        {"type": "string"},
                        {"type": "number"},
                        {"type": "boolean"},
                        {"type": "array", "items": {"type": "string"}}
                    ]
                }
            }
        },
        "inventory": {
            "type": "object",
            "properties": {
                "quantity": {"type": "integer", "minimum": 0},
                "warehouse": {"type": "string"},
                "reserved": {"type": "integer", "minimum": 0, "default": 0}
            },
            "required": ["quantity", "warehouse"]
        }
    },
    "required": ["sku", "name", "price", "category"]
}
```

### API Response Schema

```json
{
    "type": "object",
    "properties": {
        "success": {"type": "boolean"},
        "data": {
            "oneOf": [
                {"type": "object"},
                {"type": "array"},
                {"type": "null"}
            ]
        },
        "error": {
            "type": "object",
            "properties": {
                "code": {"type": "string"},
                "message": {"type": "string"},
                "details": {"type": "object"},
                "fields": {
                    "type": "object",
                    "patternProperties": {
                        ".*": {"type": "string"}
                    }
                }
            },
            "required": ["code", "message"]
        },
        "meta": {
            "type": "object",
            "properties": {
                "page": {"type": "integer", "minimum": 1},
                "limit": {"type": "integer", "minimum": 1, "maximum": 100},
                "total": {"type": "integer", "minimum": 0},
                "timestamp": {"type": "string", "format": "date-time"}
            }
        }
    },
    "required": ["success"],
    "if": {"properties": {"success": {"const": false}}},
    "then": {"required": ["error"]},
    "else": {"required": ["data"]}
}
```

## API Reference

### XeipuuvValidator

#### Methods

- `NewXeipuuvValidator(logger *logger.Logger) *XeipuuvValidator` - Create new validator
- `AddSchemaJSON(id string, schemaJSON string) error` - Add schema from JSON string
- `ValidateJSON(schemaID string, jsonData string) error` - Validate JSON string
- `ValidateData(schemaID string, data interface{}) error` - Validate Go object
- `ValidateWithDetails(schemaID string, data interface{}) (*ValidationResult, error)` - Detailed validation
- `ListSchemas() []string` - List all schema IDs
- `GetSchema(schemaID string) (*gojsonschema.Schema, bool)` - Get compiled schema
- `RemoveSchema(schemaID string) bool` - Remove schema

### SimpleHTTPMiddleware

#### Methods

- `NewSimpleHTTPMiddleware(validator *XeipuuvValidator, logger *logger.Logger) *SimpleHTTPMiddleware` - Create middleware
- `ValidateRequest(schemaID string) func(http.Handler) http.Handler` - Request validation middleware
- `ValidationHandler(schemaID string) http.HandlerFunc` - Standalone validation endpoint

## Testing

Run the test suite:

```bash
go test ./pkg/sdl/jsonschema/... -v
```

Run benchmarks:

```bash
go test ./pkg/sdl/jsonschema/... -bench=. -benchmem
```

## Error Handling

The validator returns structured errors that include:

- Error type and code
- Human-readable messages
- Validation details (field locations, values)
- Nested error information

All validation errors implement the `error` interface and can be type-asserted to `*errors.AppError` for additional information.

## Performance Considerations

- Schemas are compiled once and reused for multiple validations
- Use `ValidateData()` for Go objects when possible (avoids JSON marshaling)
- Consider schema caching strategies for high-throughput applications
- Benchmark your specific use cases

## Limitations

- Uses `xeipuuv/gojsonschema` which supports JSON Schema Draft 4
- Some advanced JSON Schema features may not be supported
- Large schemas or deeply nested objects may impact performance

## Migration Guide

If migrating from other JSON Schema libraries:

1. Update imports to use `XeipuuvValidator`
2. Replace schema compilation calls with `AddSchemaJSON()`
3. Update validation calls to use new method signatures
4. Update error handling for new error types