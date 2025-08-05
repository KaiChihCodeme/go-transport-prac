package avro

import (
	"os"
	"testing"
	"time"
)

func TestAvroManagerCreation(t *testing.T) {
	manager, err := NewManager("tmp/test_avro")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_avro")

	if manager == nil {
		t.Fatal("Manager is nil")
	}

	// Verify schemas are loaded
	userSchema := manager.GetUserSchema()
	if userSchema == nil {
		t.Fatal("User schema is nil")
	}

	productSchema := manager.GetProductSchema()
	if productSchema == nil {
		t.Fatal("Product schema is nil")
	}

	orderSchema := manager.GetOrderSchema()
	if orderSchema == nil {
		t.Fatal("Order schema is nil")
	}

	t.Log("✓ Manager created successfully with all schemas loaded")
}

func TestUserJSONSerialization(t *testing.T) {
	manager, err := NewManager("tmp/test_user_json")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_user_json")

	// Create test user
	phone := "+1-555-0123"
	user := User{
		ID:     1,
		Email:  "test@example.com",
		Name:   "Test User",
		Status: UserStatusActive,
		Profile: &Profile{
			FirstName: "Test",
			LastName:  "User",
			Phone:     &phone,
			Address: &Address{
				Street:     "123 Test St",
				City:       "Test City",
				State:      "TS",
				PostalCode: "12345",
				Country:    "USA",
			},
			Interests: []string{"testing", "avro"},
			Metadata: map[string]string{
				"test": "value",
				"env":  "test",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test JSON serialization
	jsonData, err := manager.SerializeUserJSON(user)
	if err != nil {
		t.Fatalf("Failed to serialize user to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Fatal("JSON data is empty")
	}

	t.Logf("✓ User serialized to JSON (%d bytes)", len(jsonData))

	// Test JSON deserialization
	deserializedUser, err := manager.DeserializeUserJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize user from JSON: %v", err)
	}

	// Verify data
	if deserializedUser.ID != user.ID {
		t.Errorf("ID mismatch: expected %d, got %d", user.ID, deserializedUser.ID)
	}

	if deserializedUser.Email != user.Email {
		t.Errorf("Email mismatch: expected %s, got %s", user.Email, deserializedUser.Email)
	}

	if deserializedUser.Name != user.Name {
		t.Errorf("Name mismatch: expected %s, got %s", user.Name, deserializedUser.Name)
	}

	if deserializedUser.Status != user.Status {
		t.Errorf("Status mismatch: expected %s, got %s", user.Status, deserializedUser.Status)
	}

	// Verify profile
	if deserializedUser.Profile == nil {
		t.Fatal("Profile is nil after deserialization")
	}

	if deserializedUser.Profile.FirstName != user.Profile.FirstName {
		t.Errorf("FirstName mismatch: expected %s, got %s", 
			user.Profile.FirstName, deserializedUser.Profile.FirstName)
	}

	if deserializedUser.Profile.Phone == nil {
		t.Errorf("Phone is nil after deserialization, expected: %s", *user.Profile.Phone)
	} else if *deserializedUser.Profile.Phone != *user.Profile.Phone {
		t.Errorf("Phone mismatch: expected %s, got %s", *user.Profile.Phone, *deserializedUser.Profile.Phone)
	}

	if len(deserializedUser.Profile.Interests) != len(user.Profile.Interests) {
		t.Errorf("Interests length mismatch: expected %d, got %d",
			len(user.Profile.Interests), len(deserializedUser.Profile.Interests))
	}

	if len(deserializedUser.Profile.Metadata) != len(user.Profile.Metadata) {
		t.Errorf("Metadata length mismatch: expected %d, got %d",
			len(user.Profile.Metadata), len(deserializedUser.Profile.Metadata))
	}

	t.Log("✓ User JSON serialization/deserialization successful")
}

func TestUserBinarySerialization(t *testing.T) {
	manager, err := NewManager("tmp/test_user_binary")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_user_binary")

	// Create test user  
	user := User{
		ID:     2,
		Email:  "binary@example.com",
		Name:   "Binary User",
		Status: UserStatusActive,
		Profile: &Profile{
			FirstName: "Binary",
			LastName:  "User",
			Interests: []string{"binary", "encoding"},
			Metadata: map[string]string{
				"format": "binary",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test binary serialization
	binaryData, err := manager.SerializeUserBinary(user)
	if err != nil {
		t.Fatalf("Failed to serialize user to binary: %v", err)
	}

	if len(binaryData) == 0 {
		t.Fatal("Binary data is empty")
	}

	t.Logf("✓ User serialized to binary (%d bytes)", len(binaryData))

	// Test binary deserialization
	deserializedUser, err := manager.DeserializeUserBinary(binaryData)
	if err != nil {
		t.Fatalf("Failed to deserialize user from binary: %v", err)
	}

	// Verify basic data
	if deserializedUser.ID != user.ID {
		t.Errorf("ID mismatch: expected %d, got %d", user.ID, deserializedUser.ID)
	}

	if deserializedUser.Email != user.Email {
		t.Errorf("Email mismatch: expected %s, got %s", user.Email, deserializedUser.Email)
	}

	if deserializedUser.Name != user.Name {
		t.Errorf("Name mismatch: expected %s, got %s", user.Name, deserializedUser.Name)
	}

	// Compare binary vs JSON size
	jsonData, err := manager.SerializeUserJSON(user)
	if err != nil {
		t.Fatalf("Failed to serialize user to JSON for comparison: %v", err)
	}

	t.Logf("✓ Size comparison: Binary %d bytes, JSON %d bytes", len(binaryData), len(jsonData))
	if len(binaryData) < len(jsonData) {
		savings := float64(len(jsonData)-len(binaryData)) / float64(len(jsonData)) * 100
		t.Logf("✓ Binary format saves %.1f%% space", savings)
	}

	t.Log("✓ User binary serialization/deserialization successful")
}

func TestProductSerialization(t *testing.T) {
	manager, err := NewManager("tmp/test_product")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_product")

	// Create test product
	discount := float32(10.5)
	product := Product{
		ID:          100,
		Name:        "Test Product",
		Description: "A product for testing",
		SKU:         "TEST-001",
		Price: Price{
			Currency:           "USD",
			AmountCents:        1999,
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
		Categories: []string{"Test", "Sample"},
		Tags:       []string{"test", "avro"},
		Status:     ProductStatusActive,
		Specifications: map[string]string{
			"weight": "1kg",
			"color":  "blue",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test JSON serialization
	jsonData, err := manager.SerializeProductJSON(product)
	if err != nil {
		t.Fatalf("Failed to serialize product to JSON: %v", err)
	}

	deserializedProduct, err := manager.DeserializeProductJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize product from JSON: %v", err)
	}

	// Verify basic data
	if deserializedProduct.ID != product.ID {
		t.Errorf("Product ID mismatch: expected %d, got %d", product.ID, deserializedProduct.ID)
	}

	if deserializedProduct.Name != product.Name {
		t.Errorf("Product name mismatch: expected %s, got %s", product.Name, deserializedProduct.Name)
	}

	if deserializedProduct.Price.AmountCents != product.Price.AmountCents {
		t.Errorf("Price mismatch: expected %d, got %d", 
			product.Price.AmountCents, deserializedProduct.Price.AmountCents)
	}

	// Test binary serialization
	binaryData, err := manager.SerializeProductBinary(product)
	if err != nil {
		t.Fatalf("Failed to serialize product to binary: %v", err)
	}

	deserializedProductBinary, err := manager.DeserializeProductBinary(binaryData)
	if err != nil {
		t.Fatalf("Failed to deserialize product from binary: %v", err)
	}

	if deserializedProductBinary.ID != product.ID {
		t.Errorf("Binary product ID mismatch: expected %d, got %d", 
			product.ID, deserializedProductBinary.ID)
	}

	t.Logf("✓ Product serialization: JSON %d bytes, Binary %d bytes", len(jsonData), len(binaryData))
	t.Log("✓ Product serialization/deserialization successful")
}

func TestFileOperations(t *testing.T) {
	manager, err := NewManager("tmp/test_file_ops")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_file_ops")

	// Create sample users
	users := manager.CreateSampleUsers(3)
	if len(users) != 3 {
		t.Fatalf("Expected 3 users, got %d", len(users))
	}

	// Write to file
	filename := "test_users.avro"
	err = manager.WriteUsersToFile(filename, users)
	if err != nil {
		t.Fatalf("Failed to write users to file: %v", err)
	}

	// Read from file
	readUsers, err := manager.ReadUsersFromFile(filename)
	if err != nil {
		t.Fatalf("Failed to read users from file: %v", err)
	}

	if len(readUsers) != len(users) {
		t.Fatalf("User count mismatch: expected %d, got %d", len(users), len(readUsers))
	}

	// Verify first user
	if readUsers[0].ID != users[0].ID {
		t.Errorf("First user ID mismatch: expected %d, got %d", users[0].ID, readUsers[0].ID)
	}

	if readUsers[0].Email != users[0].Email {
		t.Errorf("First user email mismatch: expected %s, got %s", users[0].Email, readUsers[0].Email)
	}

	// Test file listing
	files, err := manager.ListFiles()
	if err != nil {
		t.Fatalf("Failed to list files: %v", err)
	}

	found := false
	for _, file := range files {
		if file == filename {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find %s in file list %v", filename, files)
	}

	t.Logf("✓ File operations successful: wrote and read %d users", len(users))
}

func TestSampleDataGeneration(t *testing.T) {
	manager, err := NewManager("tmp/test_samples")
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer os.RemoveAll("tmp/test_samples")

	// Test user generation
	users := manager.CreateSampleUsers(5)
	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}

	// Verify users have required fields
	for i, user := range users {
		if user.ID == 0 {
			t.Errorf("User %d has zero ID", i)
		}
		if user.Email == "" {
			t.Errorf("User %d has empty email", i)
		}
		if user.Name == "" {
			t.Errorf("User %d has empty name", i)
		}
		if user.Status == "" {
			t.Errorf("User %d has empty status", i)
		}
		if user.Profile == nil {
			t.Errorf("User %d has nil profile", i)
		}
	}

	// Test product generation
	products := manager.CreateSampleProducts(3)
	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}

	// Verify products have required fields
	for i, product := range products {
		if product.ID == 0 {
			t.Errorf("Product %d has zero ID", i)
		}
		if product.Name == "" {
			t.Errorf("Product %d has empty name", i)
		}
		if product.SKU == "" {
			t.Errorf("Product %d has empty SKU", i)
		}
		if product.Price.Currency == "" {
			t.Errorf("Product %d has empty currency", i)
		}
		if product.Price.AmountCents == 0 {
			t.Errorf("Product %d has zero price", i)
		}
	}

	t.Log("✓ Sample data generation successful")
}