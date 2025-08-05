package parquet

import (
	"os"
	"testing"
	"time"
)

func TestSimpleParquetOperations(t *testing.T) {
	// Create test directory
	testDir := "tmp/test_simple_parquet"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	// Create sample users
	users := []User{
		{
			ID:     1,
			Email:  "test1@example.com",
			Name:   "Test User 1",
			Status: "active",
			Profile: &Profile{
				FirstName: "Test",
				LastName:  "User1",
				Phone:     "+1-555-0001",
				Address: &Address{
					Street:     "123 Test St",
					City:       "Test City",
					State:      "TS",
					PostalCode: "12345",
					Country:    "USA",
				},
				Interests: []string{"testing", "parquet"},
				Metadata: map[string]string{
					"test": "value",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:     2,
			Email:  "test2@example.com",
			Name:   "Test User 2",
			Status: "inactive",
			Profile: &Profile{
				FirstName: "Test",
				LastName:  "User2",
				Phone:     "+1-555-0002",
				Address: &Address{
					Street:     "456 Test Ave",
					City:       "Test City",
					State:      "TS",
					PostalCode: "12346",
					Country:    "USA",
				},
				Interests: []string{"data", "analytics"},
				Metadata: map[string]string{
					"department": "engineering",
				},
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	// Test write
	filename := "simple_test_users.parquet"
	err := manager.WriteUsers(filename, users)
	if err != nil {
		t.Fatalf("Failed to write users: %v", err)
	}
	t.Logf("✓ Successfully wrote %d users to %s", len(users), filename)

	// Test read
	readUsers, err := manager.ReadUsers(filename)
	if err != nil {
		t.Fatalf("Failed to read users: %v", err)
	}
	t.Logf("✓ Successfully read %d users from %s", len(readUsers), filename)

	// Verify data integrity
	if len(readUsers) != len(users) {
		t.Fatalf("Expected %d users, got %d", len(users), len(readUsers))
	}

	for i, user := range readUsers {
		original := users[i]
		
		if user.ID != original.ID {
			t.Errorf("User %d: Expected ID %d, got %d", i, original.ID, user.ID)
		}
		
		if user.Email != original.Email {
			t.Errorf("User %d: Expected email %s, got %s", i, original.Email, user.Email)
		}
		
		if user.Name != original.Name {
			t.Errorf("User %d: Expected name %s, got %s", i, original.Name, user.Name)
		}
		
		if user.Status != original.Status {
			t.Errorf("User %d: Expected status %s, got %s", i, original.Status, user.Status)
		}

		// Verify nested profile data
		if user.Profile == nil || original.Profile == nil {
			t.Errorf("User %d: Profile data missing", i)
			continue
		}
		
		if user.Profile.FirstName != original.Profile.FirstName {
			t.Errorf("User %d: Expected FirstName %s, got %s", i, original.Profile.FirstName, user.Profile.FirstName)
		}
		
		if len(user.Profile.Interests) != len(original.Profile.Interests) {
			t.Errorf("User %d: Expected %d interests, got %d", i, len(original.Profile.Interests), len(user.Profile.Interests))
		}
		
		if len(user.Profile.Metadata) != len(original.Profile.Metadata) {
			t.Errorf("User %d: Expected %d metadata entries, got %d", i, len(original.Profile.Metadata), len(user.Profile.Metadata))
		}
	}

	// Test file info
	info, err := manager.GetBasicFileInfo(filename)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.NumRows != int64(len(users)) {
		t.Errorf("Expected %d rows in file info, got %d", len(users), info.NumRows)
	}

	if info.FileSize <= 0 {
		t.Errorf("Expected positive file size, got %d", info.FileSize)
	}

	t.Logf("✓ File info: %d rows, %d bytes, %d schema fields", 
		info.NumRows, info.FileSize, len(info.Schema.Fields()))

	// Test list files
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

	t.Logf("✓ Found %d parquet files: %v", len(files), files)
	t.Logf("✓ All simple Parquet operations completed successfully")
}

func TestProductOperations(t *testing.T) {
	testDir := "tmp/test_products_parquet"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	// Create sample products
	products := []Product{
		{
			ID:          1,
			Name:        "Test Product 1", 
			Description: "A product for testing",
			SKU:         "TEST-001",
			Price: &Price{
				Currency:    "USD",
				AmountCents: 1999, // $19.99
			},
			Inventory: &Inventory{
				Quantity:       100,
				Reserved:       5,
				Available:      95,
				TrackInventory: true,
				ReorderLevel:   20,
				MaxStock:       500,
			},
			Categories: []string{"Test", "Sample"},
			Tags:       []string{"testing", "parquet"},
			Status:     "active",
			Specifications: map[string]string{
				"weight": "1.5kg",
				"color":  "blue",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	filename := "test_products.parquet"
	
	// Test write products
	err := manager.WriteProducts(filename, products)
	if err != nil {
		t.Fatalf("Failed to write products: %v", err)
	}
	t.Logf("✓ Successfully wrote %d products", len(products))

	// Test read products
	readProducts, err := manager.ReadProducts(filename)
	if err != nil {
		t.Fatalf("Failed to read products: %v", err)
	}
	t.Logf("✓ Successfully read %d products", len(readProducts))

	// Verify data
	if len(readProducts) != len(products) {
		t.Fatalf("Expected %d products, got %d", len(products), len(readProducts))
	}

	product := readProducts[0]
	original := products[0]

	if product.ID != original.ID || product.Name != original.Name {
		t.Errorf("Product data mismatch: ID %d->%d, Name %s->%s", 
			original.ID, product.ID, original.Name, product.Name)
	}

	if product.Price.AmountCents != original.Price.AmountCents {
		t.Errorf("Price mismatch: %d->%d", original.Price.AmountCents, product.Price.AmountCents)
	}

	t.Logf("✓ Product operations completed successfully")
}