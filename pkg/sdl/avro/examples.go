package avro

import (
	"fmt"
	"log"
)

// Examples demonstrates various Avro operations
type Examples struct {
	manager *Manager
}

// NewExamples creates a new examples instance
func NewExamples() (*Examples, error) {
	manager, err := NewManager("tmp/avro_examples")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	return &Examples{
		manager: manager,
	}, nil
}

// RunAllExamples runs all demonstration examples
func (e *Examples) RunAllExamples() error {
	fmt.Println("=== Avro Examples ===")

	if err := e.JSONEncodingExample(); err != nil {
		return fmt.Errorf("JSON encoding example failed: %w", err)
	}

	if err := e.BinaryEncodingExample(); err != nil {
		return fmt.Errorf("binary encoding example failed: %w", err)
	}

	if err := e.FileOperationsExample(); err != nil {
		return fmt.Errorf("file operations example failed: %w", err)
	}

	if err := e.SchemaIntrospectionExample(); err != nil {
		return fmt.Errorf("schema introspection example failed: %w", err)
	}

	if err := e.DataValidationExample(); err != nil {
		return fmt.Errorf("data validation example failed: %w", err)
	}

	if err := e.SchemaEvolutionExample(); err != nil {
		return fmt.Errorf("schema evolution example failed: %w", err)
	}

	if err := e.SchemaRegistryExample(); err != nil {
		return fmt.Errorf("schema registry example failed: %w", err)
	}

	if err := e.PerformanceComparisonExample(); err != nil {
		return fmt.Errorf("performance comparison example failed: %w", err)
	}

	fmt.Println("✓ All Avro examples completed successfully")
	return nil
}

// JSONEncodingExample demonstrates JSON encoding with Avro schema
func (e *Examples) JSONEncodingExample() error {
	fmt.Println("--- JSON Encoding Example ---")

	// Create sample user
	user := User{
		ID:     1,
		Email:  "alice@example.com",
		Name:   "Alice Johnson",
		Status: UserStatusActive,
		Profile: &Profile{
			FirstName: "Alice",
			LastName:  "Johnson",
			Phone:     stringPtr("+1-555-0123"),
			Address: &Address{
				Street:     "123 Main St",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			Interests: []string{"technology", "reading", "travel"},
			Metadata: map[string]string{
				"source":      "registration",
				"preferences": "email",
				"tier":        "premium",
			},
		},
	}

	// Serialize to JSON
	jsonData, err := e.manager.SerializeUserJSON(user)
	if err != nil {
		return fmt.Errorf("failed to serialize user to JSON: %w", err)
	}

	fmt.Printf("✓ Serialized user to JSON (%d bytes)\n", len(jsonData))
	fmt.Printf("JSON data: %s\n", string(jsonData))

	// Deserialize from JSON
	deserializedUser, err := e.manager.DeserializeUserJSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to deserialize user from JSON: %w", err)
	}

	fmt.Printf("✓ Deserialized user from JSON\n")
	fmt.Printf("  User: ID=%d, Email=%s, Name=%s\n", 
		deserializedUser.ID, deserializedUser.Email, deserializedUser.Name)
	fmt.Printf("  Profile: %s %s, Phone=%s\n",
		deserializedUser.Profile.FirstName, deserializedUser.Profile.LastName,
		*deserializedUser.Profile.Phone)
	fmt.Printf("  Address: %s, %s, %s\n",
		deserializedUser.Profile.Address.City, 
		deserializedUser.Profile.Address.State,
		deserializedUser.Profile.Address.Country)
	fmt.Printf("  Interests: %v\n", deserializedUser.Profile.Interests)
	fmt.Printf("  Metadata: %v\n", deserializedUser.Profile.Metadata)

	// Verify data integrity
	if err := e.verifyUserData(user, deserializedUser); err != nil {
		return fmt.Errorf("data integrity check failed: %w", err)
	}

	fmt.Println("✓ JSON encoding/decoding data integrity verified")
	return nil
}

// BinaryEncodingExample demonstrates binary encoding with Avro
func (e *Examples) BinaryEncodingExample() error {
	fmt.Println("--- Binary Encoding Example ---")

	// Create sample product
	discount := float32(15.0)
	product := Product{
		ID:          100,
		Name:        "Wireless Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		SKU:         "WH-001",
		Price: Price{
			Currency:           "USD",
			AmountCents:        29999, // $299.99
			DiscountPercentage: &discount,
		},
		Inventory: Inventory{
			Quantity:       50,
			Reserved:       5,
			Available:      45,
			TrackInventory: true,
			ReorderLevel:   10,
			MaxStock:       100,
		},
		Categories: []string{"Electronics", "Audio", "Headphones"},
		Tags:       []string{"wireless", "bluetooth", "noise-canceling"},
		Status:     ProductStatusActive,
		Specifications: map[string]string{
			"battery_life": "20 hours",
			"weight":       "250g",
			"color":        "black",
			"connectivity": "Bluetooth 5.0",
		},
	}

	// Serialize to binary
	binaryData, err := e.manager.SerializeProductBinary(product)
	if err != nil {
		return fmt.Errorf("failed to serialize product to binary: %w", err)
	}

	fmt.Printf("✓ Serialized product to binary (%d bytes)\n", len(binaryData))

	// Also serialize to JSON for size comparison
	jsonData, err := e.manager.SerializeProductJSON(product)
	if err != nil {
		return fmt.Errorf("failed to serialize product to JSON: %w", err)
	}

	fmt.Printf("✓ Binary size: %d bytes, JSON size: %d bytes\n", len(binaryData), len(jsonData))
	fmt.Printf("✓ Binary is %.1f%% smaller than JSON\n", 
		float64(len(jsonData)-len(binaryData))/float64(len(jsonData))*100)

	// Deserialize from binary
	deserializedProduct, err := e.manager.DeserializeProductBinary(binaryData)
	if err != nil {
		return fmt.Errorf("failed to deserialize product from binary: %w", err)
	}

	fmt.Printf("✓ Deserialized product from binary\n")
	fmt.Printf("  Product: ID=%d, Name=%s, SKU=%s\n",
		deserializedProduct.ID, deserializedProduct.Name, deserializedProduct.SKU)
	fmt.Printf("  Price: %s %.2f", deserializedProduct.Price.Currency, 
		float64(deserializedProduct.Price.AmountCents)/100)
	if deserializedProduct.Price.DiscountPercentage != nil {
		fmt.Printf(" (%.1f%% discount)", *deserializedProduct.Price.DiscountPercentage)
	}
	fmt.Printf("\n")
	fmt.Printf("  Inventory: %d available, %d reserved\n",
		deserializedProduct.Inventory.Available, deserializedProduct.Inventory.Reserved)
	fmt.Printf("  Categories: %v\n", deserializedProduct.Categories)
	fmt.Printf("  Specifications: %v\n", deserializedProduct.Specifications)

	// Verify data integrity
	if err := e.verifyProductData(product, deserializedProduct); err != nil {
		return fmt.Errorf("data integrity check failed: %w", err)
	}

	fmt.Println("✓ Binary encoding/decoding data integrity verified")
	return nil
}

// FileOperationsExample demonstrates file I/O operations
func (e *Examples) FileOperationsExample() error {
	fmt.Println("--- File Operations Example ---")

	// Create sample users
	users := e.manager.CreateSampleUsers(5)
	fmt.Printf("Created %d sample users\n", len(users))

	// Write to file
	filename := "sample_users.avro"
	err := e.manager.WriteUsersToFile(filename, users)
	if err != nil {
		return fmt.Errorf("failed to write users to file: %w", err)
	}
	fmt.Printf("✓ Wrote %d users to %s\n", len(users), filename)

	// Read from file
	readUsers, err := e.manager.ReadUsersFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read users from file: %w", err)
	}
	fmt.Printf("✓ Read %d users from %s\n", len(readUsers), filename)

	// Verify data
	if len(users) != len(readUsers) {
		return fmt.Errorf("user count mismatch: wrote %d, read %d", len(users), len(readUsers))
	}

	for i, original := range users {
		read := readUsers[i]
		if err := e.verifyUserData(original, read); err != nil {
			return fmt.Errorf("user %d data mismatch: %w", i, err)
		}
	}

	fmt.Println("✓ File operations data integrity verified")
	return nil
}

// SchemaIntrospectionExample demonstrates schema introspection
func (e *Examples) SchemaIntrospectionExample() error {
	fmt.Println("--- Schema Introspection Example ---")

	// Get user schema
	userSchema := e.manager.GetUserSchema()
	fmt.Printf("User schema type: %s\n", userSchema.Type())
	fmt.Printf("User schema: %s\n", userSchema.String())

	// Get product schema
	productSchema := e.manager.GetProductSchema()
	fmt.Printf("Product schema type: %s\n", productSchema.Type())

	// Get order schema
	orderSchema := e.manager.GetOrderSchema()
	fmt.Printf("Order schema type: %s\n", orderSchema.Type())

	fmt.Println("✓ Schema introspection completed")
	return nil
}

// DataValidationExample demonstrates data validation with schema
func (e *Examples) DataValidationExample() error {
	fmt.Println("--- Data Validation Example ---")

	// Test with valid data
	validUser := User{
		ID:     1,
		Email:  "valid@example.com",
		Name:   "Valid User",
		Status: UserStatusActive,
		Profile: &Profile{
			FirstName: "Valid",
			LastName:  "User",
			Interests: []string{"validation"},
			Metadata:  map[string]string{"test": "true"},
		},
	}

	// This should work fine
	_, err := e.manager.SerializeUserJSON(validUser)
	if err != nil {
		return fmt.Errorf("valid data serialization failed: %w", err)
	}
	fmt.Println("✓ Valid data serialized successfully")

	// Test with invalid enum (this will be caught by Go type system, not Avro)
	// Note: invalidUser.Status = "INVALID_STATUS" would cause a compile error in Go

	// Since Go's type system prevents invalid enum values, 
	// we'll simulate by creating the data directly
	fmt.Println("✓ Data validation relies on Go's type system for compile-time checks")
	fmt.Println("✓ Avro provides runtime schema validation for serialized data")

	return nil
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

func (e *Examples) verifyUserData(original, deserialized User) error {
	if original.ID != deserialized.ID {
		return fmt.Errorf("ID mismatch: %d != %d", original.ID, deserialized.ID)
	}
	if original.Email != deserialized.Email {
		return fmt.Errorf("email mismatch: %s != %s", original.Email, deserialized.Email)
	}
	if original.Name != deserialized.Name {
		return fmt.Errorf("name mismatch: %s != %s", original.Name, deserialized.Name)
	}
	if original.Status != deserialized.Status {
		return fmt.Errorf("status mismatch: %s != %s", original.Status, deserialized.Status)
	}

	// Verify profile if both exist
	if original.Profile != nil && deserialized.Profile != nil {
		if original.Profile.FirstName != deserialized.Profile.FirstName {
			return fmt.Errorf("firstName mismatch: %s != %s", 
				original.Profile.FirstName, deserialized.Profile.FirstName)
		}
		if original.Profile.LastName != deserialized.Profile.LastName {
			return fmt.Errorf("lastName mismatch: %s != %s", 
				original.Profile.LastName, deserialized.Profile.LastName)
		}

		// Check phone
		if (original.Profile.Phone == nil) != (deserialized.Profile.Phone == nil) {
			return fmt.Errorf("phone nullability mismatch")
		}
		if original.Profile.Phone != nil && deserialized.Profile.Phone != nil {
			if *original.Profile.Phone != *deserialized.Profile.Phone {
				return fmt.Errorf("phone mismatch: %s != %s", 
					*original.Profile.Phone, *deserialized.Profile.Phone)
			}
		}

		// Check interests length
		if len(original.Profile.Interests) != len(deserialized.Profile.Interests) {
			return fmt.Errorf("interests length mismatch: %d != %d",
				len(original.Profile.Interests), len(deserialized.Profile.Interests))
		}

		// Check metadata length
		if len(original.Profile.Metadata) != len(deserialized.Profile.Metadata) {
			return fmt.Errorf("metadata length mismatch: %d != %d",
				len(original.Profile.Metadata), len(deserialized.Profile.Metadata))
		}
	} else if (original.Profile == nil) != (deserialized.Profile == nil) {
		return fmt.Errorf("profile nullability mismatch")
	}

	return nil
}

func (e *Examples) verifyProductData(original, deserialized Product) error {
	if original.ID != deserialized.ID {
		return fmt.Errorf("ID mismatch: %d != %d", original.ID, deserialized.ID)
	}
	if original.Name != deserialized.Name {
		return fmt.Errorf("name mismatch: %s != %s", original.Name, deserialized.Name)
	}
	if original.SKU != deserialized.SKU {
		return fmt.Errorf("SKU mismatch: %s != %s", original.SKU, deserialized.SKU)
	}
	if original.Status != deserialized.Status {
		return fmt.Errorf("status mismatch: %s != %s", original.Status, deserialized.Status)
	}

	// Check price
	if original.Price.Currency != deserialized.Price.Currency {
		return fmt.Errorf("currency mismatch: %s != %s", 
			original.Price.Currency, deserialized.Price.Currency)
	}
	if original.Price.AmountCents != deserialized.Price.AmountCents {
		return fmt.Errorf("amount mismatch: %d != %d", 
			original.Price.AmountCents, deserialized.Price.AmountCents)
	}

	// Check discount
	if (original.Price.DiscountPercentage == nil) != (deserialized.Price.DiscountPercentage == nil) {
		return fmt.Errorf("discount nullability mismatch")
	}
	if original.Price.DiscountPercentage != nil && deserialized.Price.DiscountPercentage != nil {
		if *original.Price.DiscountPercentage != *deserialized.Price.DiscountPercentage {
			return fmt.Errorf("discount mismatch: %.2f != %.2f",
				*original.Price.DiscountPercentage, *deserialized.Price.DiscountPercentage)
		}
	}

	// Check arrays length
	if len(original.Categories) != len(deserialized.Categories) {
		return fmt.Errorf("categories length mismatch: %d != %d",
			len(original.Categories), len(deserialized.Categories))
	}
	if len(original.Tags) != len(deserialized.Tags) {
		return fmt.Errorf("tags length mismatch: %d != %d",
			len(original.Tags), len(deserialized.Tags))
	}
	if len(original.Specifications) != len(deserialized.Specifications) {
		return fmt.Errorf("specifications length mismatch: %d != %d",
			len(original.Specifications), len(deserialized.Specifications))
	}

	return nil
}

// SchemaEvolutionExample demonstrates schema evolution scenarios
func (e *Examples) SchemaEvolutionExample() error {
	fmt.Println("--- Schema Evolution Example ---")

	// Create evolution manager
	evolutionManager, err := NewEvolutionManager("tmp/avro_evolution")
	if err != nil {
		return fmt.Errorf("failed to create evolution manager: %w", err)
	}

	// Show schema versions
	evolutionManager.CompareSchemas()

	// Demonstrate evolution scenarios
	err = evolutionManager.DemonstrateSchemaEvolution()
	if err != nil {
		return fmt.Errorf("schema evolution demonstration failed: %w", err)
	}

	fmt.Println("✓ Schema evolution examples completed")
	return nil
}

// SchemaRegistryExample demonstrates schema registry concepts
func (e *Examples) SchemaRegistryExample() error {
	fmt.Println("--- Schema Registry Example ---")

	err := DemonstrateSchemaRegistry()
	if err != nil {
		return fmt.Errorf("schema registry demonstration failed: %w", err)   
	}

	fmt.Println("✓ Schema registry examples completed")
	return nil
}

// PerformanceComparisonExample demonstrates performance comparisons
func (e *Examples) PerformanceComparisonExample() error {
	fmt.Println("--- Performance Comparison Example ---")

	err := RunPerformanceComparison()
	if err != nil {
		return fmt.Errorf("performance comparison failed: %w", err)
	}

	fmt.Println("✓ Performance comparison examples completed")
	return nil
}

// CleanupExamples cleans up example files
func (e *Examples) CleanupExamples() error {
	fmt.Println("--- Cleanup Examples ---")

	files, err := e.manager.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	fmt.Printf("Cleaning up %d files...\n", len(files))

	for _, filename := range files {
		if err := e.manager.DeleteFile(filename); err != nil {
			log.Printf("Warning: failed to delete %s: %v", filename, err)
		} else {
			fmt.Printf("  Deleted %s\n", filename)
		}
	}

	fmt.Println("✓ Cleanup completed")
	return nil
}