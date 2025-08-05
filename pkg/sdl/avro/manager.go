package avro

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hamba/avro/v2"
)

// Embed schema files
//go:embed schemas/*.avsc
var schemaFiles embed.FS

// Manager handles Avro serialization and deserialization operations
type Manager struct {
	baseDir     string
	userSchema  avro.Schema
	productSchema avro.Schema
	orderSchema avro.Schema
}

// NewManager creates a new Avro manager
func NewManager(baseDir string) (*Manager, error) {
	if baseDir == "" {
		baseDir = "data/avro"
	}

	manager := &Manager{
		baseDir: baseDir,
	}

	// Load schemas
	if err := manager.loadSchemas(); err != nil {
		return nil, fmt.Errorf("failed to load schemas: %w", err)
	}

	return manager, nil
}

// loadSchemas loads all Avro schemas from embedded files
func (m *Manager) loadSchemas() error {
	// Load user schema
	userSchemaBytes, err := schemaFiles.ReadFile("schemas/user.avsc")
	if err != nil {
		return fmt.Errorf("failed to read user schema: %w", err)
	}

	m.userSchema, err = avro.Parse(string(userSchemaBytes))
	if err != nil {
		return fmt.Errorf("failed to parse user schema: %w", err)
	}

	// Load product schema
	productSchemaBytes, err := schemaFiles.ReadFile("schemas/product.avsc")
	if err != nil {
		return fmt.Errorf("failed to read product schema: %w", err)
	}

	m.productSchema, err = avro.Parse(string(productSchemaBytes))
	if err != nil {
		return fmt.Errorf("failed to parse product schema: %w", err)
	}

	// Load order schema
	orderSchemaBytes, err := schemaFiles.ReadFile("schemas/order.avsc")
	if err != nil {
		return fmt.Errorf("failed to read order schema: %w", err)
	}

	m.orderSchema, err = avro.Parse(string(orderSchemaBytes))
	if err != nil {
		return fmt.Errorf("failed to parse order schema: %w", err)
	}

	return nil
}

// ensureDir creates directory if it doesn't exist
func (m *Manager) ensureDir() error {
	return os.MkdirAll(m.baseDir, 0755)
}

// SerializeUserJSON serializes a user to JSON using Avro schema
func (m *Manager) SerializeUserJSON(user User) ([]byte, error) {
	// Convert to Avro-compatible map
	data := m.userToAvroMap(user)
	return avro.Marshal(m.userSchema, data)
}

// DeserializeUserJSON deserializes a user from JSON using Avro schema
func (m *Manager) DeserializeUserJSON(data []byte) (User, error) {
	var result interface{}
	err := avro.Unmarshal(m.userSchema, data, &result)
	if err != nil {
		return User{}, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return m.avroMapToUser(result.(map[string]interface{}))
}

// SerializeUserBinary serializes a user to binary using Avro
func (m *Manager) SerializeUserBinary(user User) ([]byte, error) {
	data := m.userToAvroMap(user)
	
	var buf bytes.Buffer
	encoder := avro.NewEncoderForSchema(m.userSchema, &buf)

	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode user: %w", err)
	}

	return buf.Bytes(), nil
}

// DeserializeUserBinary deserializes a user from binary using Avro
func (m *Manager) DeserializeUserBinary(data []byte) (User, error) {
	reader := bytes.NewReader(data)
	decoder := avro.NewDecoderForSchema(m.userSchema, reader)

	var result interface{}
	err := decoder.Decode(&result)
	if err != nil {
		return User{}, fmt.Errorf("failed to decode user: %w", err)
	}

	return m.avroMapToUser(result.(map[string]interface{}))
}

// SerializeProductJSON serializes a product to JSON using Avro schema
func (m *Manager) SerializeProductJSON(product Product) ([]byte, error) {
	data := m.productToAvroMap(product)
	return avro.Marshal(m.productSchema, data)
}

// DeserializeProductJSON deserializes a product from JSON using Avro schema
func (m *Manager) DeserializeProductJSON(data []byte) (Product, error) {
	var result interface{}
	err := avro.Unmarshal(m.productSchema, data, &result)
	if err != nil {
		return Product{}, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return m.avroMapToProduct(result.(map[string]interface{}))
}

// SerializeProductBinary serializes a product to binary using Avro
func (m *Manager) SerializeProductBinary(product Product) ([]byte, error) {
	data := m.productToAvroMap(product)
	
	var buf bytes.Buffer
	encoder := avro.NewEncoderForSchema(m.productSchema, &buf)

	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode product: %w", err)
	}

	return buf.Bytes(), nil
}

// DeserializeProductBinary deserializes a product from binary using Avro
func (m *Manager) DeserializeProductBinary(data []byte) (Product, error) {
	reader := bytes.NewReader(data)
	decoder := avro.NewDecoderForSchema(m.productSchema, reader)

	var result interface{}
	err := decoder.Decode(&result)
	if err != nil {
		return Product{}, fmt.Errorf("failed to decode product: %w", err)
	}

	return m.avroMapToProduct(result.(map[string]interface{}))
}

// WriteUsersToFile writes users to a binary Avro file
func (m *Manager) WriteUsersToFile(filename string, users []User) error {
	if err := m.ensureDir(); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := avro.NewEncoderForSchema(m.userSchema, file)

	for _, user := range users {
		data := m.userToAvroMap(user)
		err := encoder.Encode(data)
		if err != nil {
			return fmt.Errorf("failed to encode user %d: %w", user.ID, err)
		}
	}

	return nil
}

// ReadUsersFromFile reads users from a binary Avro file
func (m *Manager) ReadUsersFromFile(filename string) ([]User, error) {
	filePath := filepath.Join(m.baseDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := avro.NewDecoderForSchema(m.userSchema, file)

	var users []User
	for {
		var result interface{}
		err := decoder.Decode(&result)
		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		user, err := m.avroMapToUser(result.(map[string]interface{}))
		if err != nil {
			return nil, fmt.Errorf("failed to convert avro map to user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserSchema returns the user schema
func (m *Manager) GetUserSchema() avro.Schema {
	return m.userSchema
}

// GetProductSchema returns the product schema
func (m *Manager) GetProductSchema() avro.Schema {
	return m.productSchema
}

// GetOrderSchema returns the order schema
func (m *Manager) GetOrderSchema() avro.Schema {
	return m.orderSchema
}

// CreateSampleUsers creates sample user data for testing
func (m *Manager) CreateSampleUsers(count int) []User {
	users := make([]User, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		phone := fmt.Sprintf("+1-555-%04d", i+1000)
		users[i] = User{
			ID:     int64(i + 1),
			Email:  fmt.Sprintf("user%d@example.com", i+1),
			Name:   fmt.Sprintf("User %d", i+1),
			Status: UserStatusActive,
			Profile: &Profile{
				FirstName: fmt.Sprintf("First%d", i+1),
				LastName:  fmt.Sprintf("Last%d", i+1),
				Phone:     &phone,
				Address: &Address{
					Street:     fmt.Sprintf("%d Main St", (i+1)*100),
					City:       "Test City",
					State:      "TS",
					PostalCode: fmt.Sprintf("%05d", i+10000),
					Country:    "USA",
				},
				Interests: []string{"technology", "sports", "music"},
				Metadata: map[string]string{
					"source":    "sample_data",
					"batch_id":  fmt.Sprintf("batch_%d", i/100),
					"user_type": "standard",
				},
			},
			CreatedAt: now.Add(-time.Duration(i) * time.Hour),
			UpdatedAt: now,
		}
	}

	return users
}

// CreateSampleProducts creates sample product data for testing
func (m *Manager) CreateSampleProducts(count int) []Product {
	products := make([]Product, count)
	now := time.Now()

	categories := [][]string{
		{"Electronics", "Computers"},
		{"Clothing", "Accessories"},
		{"Books", "Fiction"},
		{"Home", "Kitchen"},
		{"Sports", "Outdoors"},
	}

	for i := 0; i < count; i++ {
		catIndex := i % len(categories)
		discountPercentage := float32(i%20) / 100.0 // 0-19%
		
		var discount *float32
		if discountPercentage > 0 {
			discount = &discountPercentage
		}

		products[i] = Product{
			ID:          int64(i + 1),
			Name:        fmt.Sprintf("Product %d", i+1),
			Description: fmt.Sprintf("Description for product %d", i+1),
			SKU:         fmt.Sprintf("SKU-%06d", i+1),
			Price: Price{
				Currency:           "USD",
				AmountCents:        int64((i%100+1) * 100), // $1.00 to $100.00
				DiscountPercentage: discount,
			},
			Inventory: Inventory{
				Quantity:       int32((i%1000) + 100),
				Reserved:       int32(i % 50),
				Available:      int32((i%1000) + 100 - (i%50)),
				TrackInventory: true,
				ReorderLevel:   int32(i%20 + 10),
				MaxStock:       int32((i%1000) + 1000),
			},
			Categories: categories[catIndex],
			Tags:       []string{"sample", "test", fmt.Sprintf("tag%d", i%10)},
			Status:     ProductStatusActive,
			Specifications: map[string]string{
				"weight": fmt.Sprintf("%.1fkg", float64(i%100+1)/10),
				"color":  []string{"red", "blue", "green", "black", "white"}[i%5],
				"size":   []string{"small", "medium", "large", "xl"}[i%4],
			},
			CreatedAt: now.Add(-time.Duration(i) * time.Hour * 24),
			UpdatedAt: now,
		}
	}

	return products
}

// ListFiles lists all Avro files in the base directory
func (m *Manager) ListFiles() ([]string, error) {
	if err := m.ensureDir(); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".avro" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// DeleteFile deletes an Avro file
func (m *Manager) DeleteFile(filename string) error {
	filePath := filepath.Join(m.baseDir, filename)
	return os.Remove(filePath)
}