package jsonschema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-transport-prac/internal/testutil"
)

func TestXeipuuvValidator_AddSchemaJSON(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	testCases := []struct {
		name       string
		schemaID   string
		schemaJSON string
		expectErr  bool
	}{
		{
			name:     "valid simple schema",
			schemaID: "user",
			schemaJSON: `{
				"type": "object",
				"properties": {
					"name": {"type": "string"},
					"age": {"type": "integer", "minimum": 0}
				},
				"required": ["name"]
			}`,
			expectErr: false,
		},
		{
			name:       "malformed JSON",
			schemaID:   "malformed",
			schemaJSON: `{"type": "object"`,
			expectErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.AddSchemaJSON(tc.schemaID, tc.schemaJSON)
			if tc.expectErr {
				helper.AssertError(err)
			} else {
				helper.AssertNoError(err)

				// Verify schema was added
				schemas := validator.ListSchemas()
				assert.Contains(t, schemas, tc.schemaID)
			}
		})
	}
}

func TestXeipuuvValidator_ValidateJSON(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	// Add a test schema
	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string", "minLength": 1},
			"age": {"type": "integer", "minimum": 0, "maximum": 150},
			"email": {"type": "string", "format": "email"}
		},
		"required": ["name", "age"]
	}`

	err := validator.AddSchemaJSON("person", schemaJSON)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		schemaID  string
		jsonData  string
		expectErr bool
	}{
		{
			name:      "valid data",
			schemaID:  "person",
			jsonData:  `{"name": "John Doe", "age": 30, "email": "john@example.com"}`,
			expectErr: false,
		},
		{
			name:      "valid data without optional field",
			schemaID:  "person",
			jsonData:  `{"name": "Jane", "age": 25}`,
			expectErr: false,
		},
		{
			name:      "missing required field",
			schemaID:  "person",
			jsonData:  `{"name": "John"}`,
			expectErr: true,
		},
		{
			name:      "invalid type",
			schemaID:  "person",
			jsonData:  `{"name": "John", "age": "thirty"}`,
			expectErr: true,
		},
		{
			name:      "value out of range",
			schemaID:  "person",
			jsonData:  `{"name": "John", "age": 200}`,
			expectErr: true,
		},
		{
			name:      "empty name",
			schemaID:  "person",
			jsonData:  `{"name": "", "age": 25}`,
			expectErr: true,
		},
		{
			name:      "invalid JSON",
			schemaID:  "person",
			jsonData:  `{"name": "John", "age": 25`,
			expectErr: true,
		},
		{
			name:      "schema not found",
			schemaID:  "nonexistent",
			jsonData:  `{"name": "John", "age": 25}`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateJSON(tc.schemaID, tc.jsonData)
			if tc.expectErr {
				helper.AssertError(err)
			} else {
				helper.AssertNoError(err)
			}
		})
	}
}

func TestXeipuuvValidator_ValidateData(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	// Add a test schema
	schemaJSON := `{
		"type": "object",
		"properties": {
			"id": {"type": "integer"},
			"active": {"type": "boolean"},
			"tags": {
				"type": "array",
				"items": {"type": "string"},
				"uniqueItems": true
			}
		},
		"required": ["id"]
	}`

	err := validator.AddSchemaJSON("item", schemaJSON)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		schemaID  string
		data      interface{}
		expectErr bool
	}{
		{
			name:     "valid data",
			schemaID: "item",
			data: map[string]interface{}{
				"id":     123,
				"active": true,
				"tags":   []string{"important", "urgent"},
			},
			expectErr: false,
		},
		{
			name:     "minimal valid data",
			schemaID: "item",
			data: map[string]interface{}{
				"id": 456,
			},
			expectErr: false,
		},
		{
			name:     "missing required field",
			schemaID: "item",
			data: map[string]interface{}{
				"active": true,
			},
			expectErr: true,
		},
		{
			name:     "wrong type",
			schemaID: "item",
			data: map[string]interface{}{
				"id":     "not-a-number",
				"active": true,
			},
			expectErr: true,
		},
		{
			name:     "duplicate array items",
			schemaID: "item",
			data: map[string]interface{}{
				"id":   789,
				"tags": []string{"tag1", "tag1"},
			},
			expectErr: true,
		},
		{
			name:     "schema not found",
			schemaID: "missing",
			data: map[string]interface{}{
				"id": 123,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateData(tc.schemaID, tc.data)
			if tc.expectErr {
				helper.AssertError(err)
			} else {
				helper.AssertNoError(err)
			}
		})
	}
}

func TestXeipuuvValidator_ValidateWithDetails(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string", "minLength": 2},
			"age": {"type": "integer", "minimum": 0}
		},
		"required": ["name"]
	}`

	err := validator.AddSchemaJSON("detailed", schemaJSON)
	require.NoError(t, err)

	// Test valid data
	validData := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	result, err := validator.ValidateWithDetails("detailed", validData)
	require.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)

	// Test invalid data
	invalidData := map[string]interface{}{
		"name": "A", // Too short
		"age":  -5,  // Below minimum
	}

	result, err = validator.ValidateWithDetails("detailed", invalidData)
	require.NoError(t, err)
	assert.False(t, result.Valid)
	assert.NotEmpty(t, result.Errors)
	assert.Len(t, result.Errors, 2) // Should have 2 validation errors
}

func TestXeipuuvValidator_ListSchemas(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	// Initially should be empty
	schemas := validator.ListSchemas()
	assert.Empty(t, schemas)

	// Add some schemas
	schemaIDs := []string{"schema1", "schema2", "schema3"}
	simpleSchema := `{"type": "object"}`

	for _, id := range schemaIDs {
		err := validator.AddSchemaJSON(id, simpleSchema)
		require.NoError(t, err)
	}

	// Should contain all added schemas
	schemas = validator.ListSchemas()
	assert.Len(t, schemas, len(schemaIDs))

	for _, id := range schemaIDs {
		assert.Contains(t, schemas, id)
	}
}

func TestXeipuuvValidator_RemoveSchema(t *testing.T) {
	helper := testutil.NewTestHelper(t)
	validator := NewXeipuuvValidator(helper.Logger())

	// Add a schema
	err := validator.AddSchemaJSON("removable", `{"type": "object"}`)
	require.NoError(t, err)

	// Verify it exists
	schemas := validator.ListSchemas()
	assert.Contains(t, schemas, "removable")

	// Remove it
	removed := validator.RemoveSchema("removable")
	assert.True(t, removed)

	// Verify it's gone
	schemas = validator.ListSchemas()
	assert.NotContains(t, schemas, "removable")

	// Try to remove non-existent schema
	removed = validator.RemoveSchema("nonexistent")
	assert.False(t, removed)
}

// Benchmark tests
func BenchmarkXeipuuvValidator_ValidateJSON(b *testing.B) {
	validator := NewXeipuuvValidator(nil)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "integer"},
			"active": {"type": "boolean"}
		},
		"required": ["name"]
	}`

	err := validator.AddSchemaJSON("benchmark", schemaJSON)
	require.NoError(b, err)

	testJSON := `{"name": "Test User", "age": 30, "active": true}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateJSON("benchmark", testJSON)
	}
}

func BenchmarkXeipuuvValidator_ValidateData(b *testing.B) {
	validator := NewXeipuuvValidator(nil)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "integer"},
			"active": {"type": "boolean"}
		},
		"required": ["name"]
	}`

	err := validator.AddSchemaJSON("benchmark", schemaJSON)
	require.NoError(b, err)

	testData := map[string]interface{}{
		"name":   "Test User",
		"age":    30,
		"active": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateData("benchmark", testData)
	}
}