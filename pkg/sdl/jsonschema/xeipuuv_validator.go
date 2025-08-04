package jsonschema

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"

	"go-transport-prac/internal/errors"
	"go-transport-prac/internal/logger"
)

// XeipuuvValidator provides JSON Schema validation using xeipuuv/gojsonschema
type XeipuuvValidator struct {
	schemas map[string]*gojsonschema.Schema
	logger  *logger.Logger
}

// NewXeipuuvValidator creates a new validator using xeipuuv/gojsonschema
func NewXeipuuvValidator(logger *logger.Logger) *XeipuuvValidator {
	return &XeipuuvValidator{
		schemas: make(map[string]*gojsonschema.Schema),
		logger:  logger,
	}
}

// AddSchemaJSON adds a schema from JSON string
func (v *XeipuuvValidator) AddSchemaJSON(id string, schemaJSON string) error {
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return errors.Wrap(err, errors.ErrorTypeValidation,
			errors.CodeValidationFailed,
			"failed to compile schema")
	}

	v.schemas[id] = schema
	return nil
}

// ValidateJSON validates a JSON string against a schema
func (v *XeipuuvValidator) ValidateJSON(schemaID string, jsonData string) error {
	schema, exists := v.schemas[schemaID]
	if !exists {
		return errors.ValidationError(errors.CodeValidationFailed,
			fmt.Sprintf("schema not found: %s", schemaID))
	}

	documentLoader := gojsonschema.NewStringLoader(jsonData)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return errors.ValidationError(errors.CodeInvalidInput,
			fmt.Sprintf("validation error: %v", err))
	}

	if !result.Valid() {
		errorMessages := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			errorMessages[i] = desc.String()
		}
		return errors.ValidationError(errors.CodeValidationFailed,
			fmt.Sprintf("validation failed: %v", errorMessages))
	}

	return nil
}

// ValidateData validates Go data against a schema
func (v *XeipuuvValidator) ValidateData(schemaID string, data interface{}) error {
	schema, exists := v.schemas[schemaID]
	if !exists {
		return errors.ValidationError(errors.CodeValidationFailed,
			fmt.Sprintf("schema not found: %s", schemaID))
	}

	documentLoader := gojsonschema.NewGoLoader(data)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return errors.ValidationError(errors.CodeInvalidInput,
			fmt.Sprintf("validation error: %v", err))
	}

	if !result.Valid() {
		errorMessages := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			errorMessages[i] = desc.String()
		}
		return errors.ValidationError(errors.CodeValidationFailed,
			fmt.Sprintf("validation failed: %v", errorMessages))
	}

	return nil
}

// ValidateWithDetails returns detailed validation results
func (v *XeipuuvValidator) ValidateWithDetails(schemaID string, data interface{}) (*ValidationResult, error) {
	schema, exists := v.schemas[schemaID]
	if !exists {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Message: fmt.Sprintf("schema not found: %s", schemaID),
				},
			},
		}, nil
	}

	documentLoader := gojsonschema.NewGoLoader(data)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return nil, errors.ValidationError(errors.CodeInvalidInput,
			fmt.Sprintf("validation error: %v", err))
	}

	validationResult := &ValidationResult{
		Valid:  result.Valid(),
		Schema: schemaID,
		Data:   data,
	}

	if !result.Valid() {
		validationResult.Errors = make([]ValidationError, len(result.Errors()))
		for i, desc := range result.Errors() {
			validationResult.Errors[i] = ValidationError{
				InstanceLocation: desc.Field(), // This maps to the JSON pointer
				Message:          desc.Description(),
				Value:            desc.Value(),
			}
		}
	}

	return validationResult, nil
}

// ListSchemas returns all registered schema IDs
func (v *XeipuuvValidator) ListSchemas() []string {
	ids := make([]string, 0, len(v.schemas))
	for id := range v.schemas {
		ids = append(ids, id)
	}
	return ids
}

// GetSchema returns a compiled schema by ID
func (v *XeipuuvValidator) GetSchema(schemaID string) (*gojsonschema.Schema, bool) {
	schema, exists := v.schemas[schemaID]
	return schema, exists
}

// RemoveSchema removes a schema from the validator
func (v *XeipuuvValidator) RemoveSchema(schemaID string) bool {
	if _, exists := v.schemas[schemaID]; exists {
		delete(v.schemas, schemaID)
		return true
	}
	return false
}

// ValidationResult represents validation results
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []ValidationError  `json:"errors,omitempty"`
	Schema string             `json:"schema,omitempty"`
	Data   interface{}        `json:"data,omitempty"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	InstanceLocation string      `json:"instance_location"`
	KeywordLocation  string      `json:"keyword_location,omitempty"`
	Message          string      `json:"message"`
	Value            interface{} `json:"value,omitempty"`
	Schema           interface{} `json:"schema,omitempty"`
}