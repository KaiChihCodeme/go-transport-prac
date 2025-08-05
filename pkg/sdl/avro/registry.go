package avro

import (
	"fmt"
	"sync"
	"time"

	"github.com/hamba/avro/v2"
)

// SchemaRegistry simulates a schema registry for managing Avro schemas
type SchemaRegistry struct {
	mu              sync.RWMutex
	schemas         map[int]SchemaMetadata
	subjectSchemas  map[string][]int
	nextSchemaID    int
	compatibilityLevels map[string]CompatibilityLevel
}

// SchemaMetadata contains metadata about a registered schema
type SchemaMetadata struct {
	ID          int                 `json:"id"`
	Version     int                 `json:"version"`
	Subject     string              `json:"subject"`
	Schema      avro.Schema         `json:"-"`
	SchemaJSON  string              `json:"schema"`
	CreatedAt   time.Time           `json:"createdAt"`
	Fingerprint string              `json:"fingerprint"`
	References  []SchemaReference   `json:"references,omitempty"`
}

// SchemaReference represents a reference to another schema
type SchemaReference struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
}

// CompatibilityLevel defines schema compatibility levels
type CompatibilityLevel string

const (
	CompatibilityNone     CompatibilityLevel = "NONE"
	CompatibilityFull     CompatibilityLevel = "FULL"
	CompatibilityForward  CompatibilityLevel = "FORWARD"
	CompatibilityBackward CompatibilityLevel = "BACKWARD"
)

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		schemas:             make(map[int]SchemaMetadata),
		subjectSchemas:     make(map[string][]int),
		nextSchemaID:       1,
		compatibilityLevels: make(map[string]CompatibilityLevel),
	}
}

// RegisterSchema registers a new schema or returns existing schema ID
func (sr *SchemaRegistry) RegisterSchema(subject string, schemaJSON string) (int, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	// Parse the schema to validate it
	schema, err := avro.Parse(schemaJSON)
	if err != nil {
		return 0, fmt.Errorf("invalid schema: %w", err)
	}

	// Generate fingerprint (simplified - in real implementation would use actual fingerprinting)
	fingerprint := fmt.Sprintf("fp_%s_%d", subject, len(schemaJSON))

	// Check if schema already exists for this subject
	if schemaIDs, exists := sr.subjectSchemas[subject]; exists {
		for _, id := range schemaIDs {
			if sr.schemas[id].Fingerprint == fingerprint {
				return id, nil // Schema already registered
			}
		}
	}

	// Check compatibility with existing schemas
	if err := sr.checkCompatibility(subject, schema); err != nil {
		return 0, fmt.Errorf("schema compatibility check failed: %w", err)
	}

	// Register new schema
	schemaID := sr.nextSchemaID
	sr.nextSchemaID++

	version := len(sr.subjectSchemas[subject]) + 1

	metadata := SchemaMetadata{
		ID:          schemaID,
		Version:     version,
		Subject:     subject,
		Schema:      schema,
		SchemaJSON:  schemaJSON,
		CreatedAt:   time.Now(),
		Fingerprint: fingerprint,
	}

	sr.schemas[schemaID] = metadata
	sr.subjectSchemas[subject] = append(sr.subjectSchemas[subject], schemaID)

	return schemaID, nil
}

// GetSchema retrieves a schema by ID
func (sr *SchemaRegistry) GetSchema(schemaID int) (SchemaMetadata, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	metadata, exists := sr.schemas[schemaID]
	if !exists {
		return SchemaMetadata{}, fmt.Errorf("schema with ID %d not found", schemaID)
	}

	return metadata, nil
}

// GetLatestSchema retrieves the latest schema for a subject
func (sr *SchemaRegistry) GetLatestSchema(subject string) (SchemaMetadata, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	schemaIDs, exists := sr.subjectSchemas[subject]
	if !exists || len(schemaIDs) == 0 {
		return SchemaMetadata{}, fmt.Errorf("no schemas found for subject %s", subject)
	}

	latestID := schemaIDs[len(schemaIDs)-1]
	return sr.schemas[latestID], nil
}

// GetSchemaVersion retrieves a specific version of a schema for a subject
func (sr *SchemaRegistry) GetSchemaVersion(subject string, version int) (SchemaMetadata, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	schemaIDs, exists := sr.subjectSchemas[subject]
	if !exists || version < 1 || version > len(schemaIDs) {
		return SchemaMetadata{}, fmt.Errorf("schema version %d not found for subject %s", version, subject)
	}

	schemaID := schemaIDs[version-1]
	return sr.schemas[schemaID], nil
}

// ListSubjects returns all registered subjects
func (sr *SchemaRegistry) ListSubjects() []string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	subjects := make([]string, 0, len(sr.subjectSchemas))
	for subject := range sr.subjectSchemas {
		subjects = append(subjects, subject)
	}
	return subjects
}

// ListSchemaVersions returns all versions for a subject
func (sr *SchemaRegistry) ListSchemaVersions(subject string) ([]int, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	schemaIDs, exists := sr.subjectSchemas[subject]
	if !exists {
		return nil, fmt.Errorf("subject %s not found", subject)
	}

	versions := make([]int, len(schemaIDs))
	for i, id := range schemaIDs {
		versions[i] = sr.schemas[id].Version
	}
	return versions, nil
}

// SetCompatibilityLevel sets the compatibility level for a subject
func (sr *SchemaRegistry) SetCompatibilityLevel(subject string, level CompatibilityLevel) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.compatibilityLevels[subject] = level
	return nil
}

// GetCompatibilityLevel gets the compatibility level for a subject
func (sr *SchemaRegistry) GetCompatibilityLevel(subject string) CompatibilityLevel {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	if level, exists := sr.compatibilityLevels[subject]; exists {
		return level
	}
	return CompatibilityBackward // Default compatibility level
}

// CheckCompatibility checks if a new schema is compatible with existing schemas
func (sr *SchemaRegistry) CheckCompatibility(subject string, schemaJSON string) (bool, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	schema, err := avro.Parse(schemaJSON)
	if err != nil {
		return false, fmt.Errorf("invalid schema: %w", err)
	}

	return sr.checkCompatibility(subject, schema) == nil, nil
}

// checkCompatibility performs the actual compatibility check
// Note: This method assumes the caller already holds the appropriate lock
func (sr *SchemaRegistry) checkCompatibility(subject string, newSchema avro.Schema) error {
	// Get compatibility level without additional locking since caller holds lock
	compatibilityLevel := CompatibilityBackward // Default
	if level, exists := sr.compatibilityLevels[subject]; exists {
		compatibilityLevel = level
	}
	
	// If no compatibility checking required
	if compatibilityLevel == CompatibilityNone {
		return nil
	}

	schemaIDs, exists := sr.subjectSchemas[subject]
	if !exists || len(schemaIDs) == 0 {
		return nil // No existing schemas to check against
	}

	// Get the latest schema for compatibility checking
	latestID := schemaIDs[len(schemaIDs)-1]
	latestSchema := sr.schemas[latestID].Schema

	switch compatibilityLevel {
	case CompatibilityForward:
		return sr.checkForwardCompatibility(latestSchema, newSchema)
	case CompatibilityBackward:
		return sr.checkBackwardCompatibility(latestSchema, newSchema)
	case CompatibilityFull:
		if err := sr.checkForwardCompatibility(latestSchema, newSchema); err != nil {
			return err
		}
		return sr.checkBackwardCompatibility(latestSchema, newSchema)
	default:
		return nil
	}
}

// checkForwardCompatibility checks if new schema can read data written with old schema
func (sr *SchemaRegistry) checkForwardCompatibility(oldSchema, newSchema avro.Schema) error {
	// Simplified compatibility check - in practice this would be more comprehensive
	if oldSchema.Type() != newSchema.Type() {
		return fmt.Errorf("schema types don't match: %s vs %s", oldSchema.Type(), newSchema.Type())
	}
	
	// Check if schemas are identical (simplified check)
	if oldSchema.String() == newSchema.String() {
		return nil
	}

	// This is a simplified check. Real implementation would:
	// - Check field additions/removals
	// - Verify default values
	// - Check type promotions
	// - Validate enum symbol additions
	
	fmt.Printf("⚠ Forward compatibility check passed (simplified)\n")
	return nil
}

// checkBackwardCompatibility checks if old schema can read data written with new schema
func (sr *SchemaRegistry) checkBackwardCompatibility(oldSchema, newSchema avro.Schema) error {
	// Simplified compatibility check
	if oldSchema.Type() != newSchema.Type() {
		return fmt.Errorf("schema types don't match: %s vs %s", oldSchema.Type(), newSchema.Type())
	}

	// Check if schemas are identical (simplified check)
	if oldSchema.String() == newSchema.String() {
		return nil
	}

	// This is a simplified check. Real implementation would:
	// - Ensure no required fields were added
	// - Check that removed fields had defaults
	// - Verify enum symbol compatibility
	// - Check type compatibility
	
	fmt.Printf("⚠ Backward compatibility check passed (simplified)\n")
	return nil
}

// GetStats returns registry statistics
func (sr *SchemaRegistry) GetStats() map[string]interface{} {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	stats := map[string]interface{}{
		"total_schemas":     len(sr.schemas),
		"total_subjects":    len(sr.subjectSchemas),
		"next_schema_id":    sr.nextSchemaID,
		"subjects":          sr.ListSubjects(),
	}

	subjectStats := make(map[string]int)
	for subject, schemaIDs := range sr.subjectSchemas {
		subjectStats[subject] = len(schemaIDs)
	}
	stats["schemas_per_subject"] = subjectStats

	return stats
}

// DemonstrateSchemaRegistry shows how a schema registry would work
func DemonstrateSchemaRegistry() error {
	fmt.Println("=== Schema Registry Demonstration ===")

	registry := NewSchemaRegistry()

	// Set compatibility levels
	registry.SetCompatibilityLevel("user", CompatibilityBackward)
	registry.SetCompatibilityLevel("product", CompatibilityFull)

	// Register schemas
	fmt.Println("--- Registering Schemas ---")

	// Read schema files
	userV1Schema, err := schemaFiles.ReadFile("schemas/user.avsc")
	if err != nil {
		return fmt.Errorf("failed to read user v1 schema: %w", err)
	}

	userV2Schema, err := evolutionSchemaFiles.ReadFile("schemas/user_v2.avsc")
	if err != nil {
		return fmt.Errorf("failed to read user v2 schema: %w", err)
	}

	productSchema, err := schemaFiles.ReadFile("schemas/product.avsc")
	if err != nil {
		return fmt.Errorf("failed to read product schema: %w", err)
	}

	// Register user v1 schema
	userV1ID, err := registry.RegisterSchema("user", string(userV1Schema))
	if err != nil {
		return fmt.Errorf("failed to register user v1 schema: %w", err)
	}
	fmt.Printf("✓ Registered user v1 schema with ID: %d\n", userV1ID)

	// Register product schema
	productID, err := registry.RegisterSchema("product", string(productSchema))
	if err != nil {
		return fmt.Errorf("failed to register product schema: %w", err)
	}
	fmt.Printf("✓ Registered product schema with ID: %d\n", productID)

	// Try to register user v2 schema (should pass backward compatibility)
	userV2ID, err := registry.RegisterSchema("user", string(userV2Schema))
	if err != nil {
		return fmt.Errorf("failed to register user v2 schema: %w", err)
	}
	fmt.Printf("✓ Registered user v2 schema with ID: %d\n", userV2ID)

	// Test schema retrieval
	fmt.Println("--- Schema Retrieval ---")

	latestUser, err := registry.GetLatestSchema("user")
	if err != nil {
		return fmt.Errorf("failed to get latest user schema: %w", err)
	}
	fmt.Printf("✓ Latest user schema: v%d (ID: %d)\n", latestUser.Version, latestUser.ID)

	userV1Retrieved, err := registry.GetSchemaVersion("user", 1)
	if err != nil {
		return fmt.Errorf("failed to get user v1 schema: %w", err)
	}
	fmt.Printf("✓ Retrieved user v1: v%d (ID: %d)\n", userV1Retrieved.Version, userV1Retrieved.ID)

	// Show registry statistics
	fmt.Println("--- Registry Statistics ---")
	stats := registry.GetStats()
	fmt.Printf("✓ Total schemas: %v\n", stats["total_schemas"])
	fmt.Printf("✓ Total subjects: %v\n", stats["total_subjects"])
	fmt.Printf("✓ Subjects: %v\n", stats["subjects"])
	fmt.Printf("✓ Schemas per subject: %v\n", stats["schemas_per_subject"])

	// Show compatibility levels
	fmt.Println("--- Compatibility Levels ---")
	fmt.Printf("✓ User compatibility: %s\n", registry.GetCompatibilityLevel("user"))
	fmt.Printf("✓ Product compatibility: %s\n", registry.GetCompatibilityLevel("product"))

	fmt.Println("✓ Schema registry demonstration completed")
	return nil
}