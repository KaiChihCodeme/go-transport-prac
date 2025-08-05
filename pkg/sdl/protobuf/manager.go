package protobuf

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-transport-prac/pkg/sdl/protobuf/gen/order"
	"go-transport-prac/pkg/sdl/protobuf/gen/product"
	"go-transport-prac/pkg/sdl/protobuf/gen/user"
)

// Manager handles Protocol Buffers serialization and deserialization
type Manager struct{}

// NewManager creates a new protobuf manager
func NewManager() *Manager {
	return &Manager{}
}

// SerializeUser serializes a User message to bytes
func (m *Manager) SerializeUser(u *user.User) ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}

	return proto.Marshal(u)
}

// DeserializeUser deserializes bytes to a User message
func (m *Manager) DeserializeUser(data []byte) (*user.User, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	u := &user.User{}
	if err := proto.Unmarshal(data, u); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return u, nil
}

// SerializeProduct serializes a Product message to bytes
func (m *Manager) SerializeProduct(p *product.Product) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("product cannot be nil")
	}

	return proto.Marshal(p)
}

// DeserializeProduct deserializes bytes to a Product message
func (m *Manager) DeserializeProduct(data []byte) (*product.Product, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	p := &product.Product{}
	if err := proto.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return p, nil
}

// SerializeOrder serializes an Order message to bytes
func (m *Manager) SerializeOrder(o *order.Order) ([]byte, error) {
	if o == nil {
		return nil, fmt.Errorf("order cannot be nil")
	}

	return proto.Marshal(o)
}

// DeserializeOrder deserializes bytes to an Order message
func (m *Manager) DeserializeOrder(data []byte) (*order.Order, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	o := &order.Order{}
	if err := proto.Unmarshal(data, o); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return o, nil
}

// Generic serialization method
func (m *Manager) Serialize(msg proto.Message) ([]byte, error) {
	if msg == nil {
		return nil, fmt.Errorf("message cannot be nil")
	}

	return proto.Marshal(msg)
}

// Generic deserialization method
func (m *Manager) Deserialize(data []byte, msg proto.Message) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	return proto.Unmarshal(data, msg)
}

// Helper functions for creating common objects

// CreateSampleUser creates a sample user for testing
func (m *Manager) CreateSampleUser() *user.User {
	now := timestamppb.Now()

	return &user.User{
		Id:     1,
		Email:  "john.doe@example.com",
		Name:   "John Doe",
		Status: user.UserStatus_USER_STATUS_ACTIVE,
		Profile: &user.Profile{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+1-555-0123",
			Address: &user.Address{
				Street:     "123 Main St",
				City:       "San Francisco",
				State:      "CA",
				PostalCode: "94105",
				Country:    "USA",
			},
			Interests: []string{"technology", "programming", "travel"},
			Metadata: map[string]string{
				"preferred_language": "en",
				"timezone":           "America/Los_Angeles",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateSampleProduct creates a sample product for testing
func (m *Manager) CreateSampleProduct() *product.Product {
	now := timestamppb.Now()

	return &product.Product{
		Id:          1,
		Name:        "Premium Wireless Headphones",
		Description: "High-quality wireless headphones with noise cancellation",
		Sku:         "WH-1000XM5",
		Price: &product.Price{
			Currency:    "USD",
			AmountCents: 39999, // $399.99
		},
		Inventory: &product.Inventory{
			Quantity:       100,
			Reserved:       5,
			Available:      95,
			TrackInventory: true,
			ReorderLevel:   20,
			MaxStock:       500,
		},
		Categories: []string{"Electronics", "Audio", "Headphones"},
		Tags:       []string{"wireless", "bluetooth", "noise-canceling", "premium"},
		Status:     product.ProductStatus_PRODUCT_STATUS_ACTIVE,
		Specifications: &product.Specifications{
			Attributes: map[string]string{
				"brand":        "Sony",
				"model":        "WH-1000XM5",
				"connectivity": "Bluetooth 5.2",
				"battery_life": "30 hours",
				"color":        "Black",
			},
			Dimensions: &product.Dimensions{
				Length: 25.4,
				Width:  22.0,
				Height: 8.5,
				Unit:   "cm",
			},
			Weight: &product.Weight{
				Value: 250,
				Unit:  "g",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateSampleOrder creates a sample order for testing
func (m *Manager) CreateSampleOrder() *order.Order {
	now := timestamppb.Now()
	deliveryTime := timestamppb.New(time.Now().Add(5 * 24 * time.Hour)) // 5 days from now

	return &order.Order{
		Id:          1,
		UserId:      1,
		OrderNumber: "ORD-2024-001234",
		Status:      order.OrderStatus_ORDER_STATUS_CONFIRMED,
		Items: []*order.OrderItem{
			{
				ProductId:   1,
				ProductName: "Premium Wireless Headphones",
				ProductSku:  "WH-1000XM5",
				Quantity:    2,
				UnitPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 39999,
				},
				TotalPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 79998,
				},
				ProductVariant: map[string]string{
					"color": "Black",
				},
			},
		},
		Summary: &order.OrderSummary{
			Subtotal: &product.Price{
				Currency:    "USD",
				AmountCents: 79998,
			},
			Tax: &product.Price{
				Currency:    "USD",
				AmountCents: 7200, // 9% tax
			},
			ShippingCost: &product.Price{
				Currency:    "USD",
				AmountCents: 999, // $9.99
			},
			Discount: &product.Price{
				Currency:    "USD",
				AmountCents: 0,
			},
			Total: &product.Price{
				Currency:    "USD",
				AmountCents: 88197, // $881.97
			},
			TotalItems: 2,
		},
		Shipping: &order.ShippingInfo{
			Address: &user.Address{
				Street:     "123 Main St",
				City:       "San Francisco",
				State:      "CA",
				PostalCode: "94105",
				Country:    "USA",
			},
			Method:         "standard",
			TrackingNumber: "1Z999AA1234567890",
			Carrier:        "ups",
			Cost: &product.Price{
				Currency:    "USD",
				AmountCents: 999,
			},
			EstimatedDelivery: deliveryTime,
		},
		Payment: &order.PaymentInfo{
			Method:        "credit_card",
			Status:        order.PaymentStatus_PAYMENT_STATUS_CAPTURED,
			TransactionId: "txn_1234567890abcdef",
			Amount: &product.Price{
				Currency:    "USD",
				AmountCents: 88197,
			},
			ProcessedAt: now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
