package protobuf

import (
	"testing"

	"go-transport-prac/pkg/sdl/protobuf/gen/user"
)

func TestManager_UserSerialization(t *testing.T) {
	manager := NewManager()
	
	// Test with sample user
	originalUser := manager.CreateSampleUser()
	
	// Serialize
	data, err := manager.SerializeUser(originalUser)
	if err != nil {
		t.Fatalf("Failed to serialize user: %v", err)
	}
	
	if len(data) == 0 {
		t.Fatal("Serialized data is empty")
	}
	
	// Deserialize
	deserializedUser, err := manager.DeserializeUser(data)
	if err != nil {
		t.Fatalf("Failed to deserialize user: %v", err)
	}
	
	// Verify key fields
	if originalUser.Id != deserializedUser.Id {
		t.Errorf("User ID mismatch: got %d, want %d", deserializedUser.Id, originalUser.Id)
	}
	
	if originalUser.Email != deserializedUser.Email {
		t.Errorf("User email mismatch: got %s, want %s", deserializedUser.Email, originalUser.Email)
	}
	
	if originalUser.Name != deserializedUser.Name {
		t.Errorf("User name mismatch: got %s, want %s", deserializedUser.Name, originalUser.Name)
	}
	
	if originalUser.Status != deserializedUser.Status {
		t.Errorf("User status mismatch: got %v, want %v", deserializedUser.Status, originalUser.Status)
	}
}

func TestManager_ProductSerialization(t *testing.T) {
	manager := NewManager()
	
	// Test with sample product
	originalProduct := manager.CreateSampleProduct()
	
	// Serialize
	data, err := manager.SerializeProduct(originalProduct)
	if err != nil {
		t.Fatalf("Failed to serialize product: %v", err)
	}
	
	if len(data) == 0 {
		t.Fatal("Serialized data is empty")
	}
	
	// Deserialize
	deserializedProduct, err := manager.DeserializeProduct(data)
	if err != nil {
		t.Fatalf("Failed to deserialize product: %v", err)
	}
	
	// Verify key fields
	if originalProduct.Id != deserializedProduct.Id {
		t.Errorf("Product ID mismatch: got %d, want %d", deserializedProduct.Id, originalProduct.Id)
	}
	
	if originalProduct.Name != deserializedProduct.Name {
		t.Errorf("Product name mismatch: got %s, want %s", deserializedProduct.Name, originalProduct.Name)
	}
	
	if originalProduct.Sku != deserializedProduct.Sku {
		t.Errorf("Product SKU mismatch: got %s, want %s", deserializedProduct.Sku, originalProduct.Sku)
	}
	
	if originalProduct.Status != deserializedProduct.Status {
		t.Errorf("Product status mismatch: got %v, want %v", deserializedProduct.Status, originalProduct.Status)
	}
}

func TestManager_NilInputs(t *testing.T) {
	manager := NewManager()
	
	// Test nil user serialization
	_, err := manager.SerializeUser(nil)
	if err == nil {
		t.Error("Expected error for nil user, got none")
	}
	
	// Test nil product serialization
	_, err = manager.SerializeProduct(nil)
	if err == nil {
		t.Error("Expected error for nil product, got none")
	}
	
	// Test nil order serialization
	_, err = manager.SerializeOrder(nil)
	if err == nil {
		t.Error("Expected error for nil order, got none")
	}
	
	// Test empty data deserialization
	_, err = manager.DeserializeUser([]byte{})
	if err == nil {
		t.Error("Expected error for empty data, got none")
	}
	
	_, err = manager.DeserializeProduct([]byte{})
	if err == nil {
		t.Error("Expected error for empty data, got none")
	}
	
	_, err = manager.DeserializeOrder([]byte{})
	if err == nil {
		t.Error("Expected error for empty data, got none")
	}
}

func TestManager_InvalidData(t *testing.T) {
	manager := NewManager()
	
	invalidData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	
	// Test invalid data deserialization
	_, err := manager.DeserializeUser(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data, got none")
	}
	
	_, err = manager.DeserializeProduct(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data, got none")
	}
	
	_, err = manager.DeserializeOrder(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data, got none")
	}
}

func TestManager_GenericSerialization(t *testing.T) {
	manager := NewManager()
	
	// Test with user
	originalUser := manager.CreateSampleUser()
	
	data, err := manager.Serialize(originalUser)
	if err != nil {
		t.Fatalf("Failed to serialize user: %v", err)
	}
	
	deserializedUser := &user.User{}
	err = manager.Deserialize(data, deserializedUser)
	if err != nil {
		t.Fatalf("Failed to deserialize user: %v", err)
	}
	
	if originalUser.Id != deserializedUser.Id {
		t.Errorf("User ID mismatch: got %d, want %d", deserializedUser.Id, originalUser.Id)
	}
}

func TestManager_SampleObjectCreation(t *testing.T) {
	manager := NewManager()
	
	// Test sample user creation
	user := manager.CreateSampleUser()
	if user == nil {
		t.Fatal("CreateSampleUser returned nil")
	}
	
	if user.Id == 0 {
		t.Error("Sample user ID should not be zero")
	}
	
	if user.Email == "" {
		t.Error("Sample user email should not be empty")
	}
	
	if user.Profile == nil {
		t.Error("Sample user profile should not be nil")
	}
	
	// Test sample product creation
	product := manager.CreateSampleProduct()
	if product == nil {
		t.Fatal("CreateSampleProduct returned nil")
	}
	
	if product.Id == 0 {
		t.Error("Sample product ID should not be zero")
	}
	
	if product.Name == "" {
		t.Error("Sample product name should not be empty")
	}
	
	if product.Price == nil {
		t.Error("Sample product price should not be nil")
	}
}

func BenchmarkUserSerialization(b *testing.B) {
	manager := NewManager()
	user := manager.CreateSampleUser()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.SerializeUser(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserDeserialization(b *testing.B) {
	manager := NewManager()
	user := manager.CreateSampleUser()
	data, err := manager.SerializeUser(user)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.DeserializeUser(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProductSerialization(b *testing.B) {
	manager := NewManager()
	product := manager.CreateSampleProduct()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.SerializeProduct(product)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProductDeserialization(b *testing.B) {
	manager := NewManager()
	product := manager.CreateSampleProduct()
	data, err := manager.SerializeProduct(product)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := manager.DeserializeProduct(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}