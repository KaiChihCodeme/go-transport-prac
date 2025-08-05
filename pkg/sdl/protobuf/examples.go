package protobuf

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"go-transport-prac/pkg/sdl/protobuf/gen/order"
	"go-transport-prac/pkg/sdl/protobuf/gen/product"
	"go-transport-prac/pkg/sdl/protobuf/gen/user"
)

// Examples demonstrates various protobuf serialization/deserialization operations
type Examples struct {
	manager *Manager
}

// NewExamples creates a new examples instance
func NewExamples() *Examples {
	return &Examples{
		manager: NewManager(),
	}
}

// RunAllExamples runs all protobuf examples
func (e *Examples) RunAllExamples() error {
	fmt.Println("=== Protocol Buffers Serialization/Deserialization Examples ===")

	if err := e.UserExample(); err != nil {
		return fmt.Errorf("user example failed: %w", err)
	}

	if err := e.ProductExample(); err != nil {
		return fmt.Errorf("product example failed: %w", err)
	}

	if err := e.OrderExample(); err != nil {
		return fmt.Errorf("order example failed: %w", err)
	}

	if err := e.SerializationSizeComparison(); err != nil {
		return fmt.Errorf("size comparison failed: %w", err)
	}

	return nil
}

// UserExample demonstrates user serialization/deserialization
func (e *Examples) UserExample() error {
	fmt.Println("--- User Example ---")

	// Create sample user
	originalUser := e.manager.CreateSampleUser()
	fmt.Printf("Original User: %+v\n", originalUser)

	// Serialize to bytes
	data, err := e.manager.SerializeUser(originalUser)
	if err != nil {
		return fmt.Errorf("failed to serialize user: %w", err)
	}
	fmt.Printf("Serialized size: %d bytes\n", len(data))

	// Deserialize back to object
	deserializedUser, err := e.manager.DeserializeUser(data)
	if err != nil {
		return fmt.Errorf("failed to deserialize user: %w", err)
	}
	fmt.Printf("Deserialized User: %+v\n", deserializedUser)

	// Verify data integrity
	if originalUser.Id != deserializedUser.Id || 
	   originalUser.Email != deserializedUser.Email ||
	   originalUser.Name != deserializedUser.Name {
		return fmt.Errorf("data integrity check failed")
	}

	fmt.Println("✓ User serialization/deserialization successful")
	return nil
}

// ProductExample demonstrates product serialization/deserialization
func (e *Examples) ProductExample() error {
	fmt.Println("--- Product Example ---")

	// Create sample product
	originalProduct := e.manager.CreateSampleProduct()
	fmt.Printf("Original Product: %+v\n", originalProduct)

	// Serialize to bytes
	data, err := e.manager.SerializeProduct(originalProduct)
	if err != nil {
		return fmt.Errorf("failed to serialize product: %w", err)
	}
	fmt.Printf("Serialized size: %d bytes\n", len(data))

	// Deserialize back to object
	deserializedProduct, err := e.manager.DeserializeProduct(data)
	if err != nil {
		return fmt.Errorf("failed to deserialize product: %w", err)
	}
	fmt.Printf("Deserialized Product: %+v\n", deserializedProduct)

	// Verify data integrity
	if originalProduct.Id != deserializedProduct.Id || 
	   originalProduct.Name != deserializedProduct.Name ||
	   originalProduct.Sku != deserializedProduct.Sku {
		return fmt.Errorf("data integrity check failed")
	}

	fmt.Println("✓ Product serialization/deserialization successful")
	return nil
}

// OrderExample demonstrates order serialization/deserialization
func (e *Examples) OrderExample() error {
	fmt.Println("--- Order Example ---")

	// Create sample order
	originalOrder := e.createSampleOrder()
	fmt.Printf("Original Order: %+v\n", originalOrder)

	// Serialize to bytes
	data, err := e.manager.SerializeOrder(originalOrder)
	if err != nil {
		return fmt.Errorf("failed to serialize order: %w", err)
	}
	fmt.Printf("Serialized size: %d bytes\n", len(data))

	// Deserialize back to object
	deserializedOrder, err := e.manager.DeserializeOrder(data)
	if err != nil {
		return fmt.Errorf("failed to deserialize order: %w", err)
	}
	fmt.Printf("Deserialized Order: %+v\n", deserializedOrder)

	// Verify data integrity
	if originalOrder.Id != deserializedOrder.Id || 
	   originalOrder.OrderNumber != deserializedOrder.OrderNumber ||
	   originalOrder.UserId != deserializedOrder.UserId {
		return fmt.Errorf("data integrity check failed")
	}

	fmt.Println("✓ Order serialization/deserialization successful")
	return nil
}

// SerializationSizeComparison compares protobuf sizes with different data types
func (e *Examples) SerializationSizeComparison() error {
	fmt.Println("--- Serialization Size Comparison ---")

	// Create test objects
	user := e.manager.CreateSampleUser()
	product := e.manager.CreateSampleProduct()
	order := e.createSampleOrder()

	// Serialize each
	userData, _ := e.manager.SerializeUser(user)
	productData, _ := e.manager.SerializeProduct(product)
	orderData, _ := e.manager.SerializeOrder(order)

	fmt.Printf("User serialized size: %d bytes\n", len(userData))
	fmt.Printf("Product serialized size: %d bytes\n", len(productData))
	fmt.Printf("Order serialized size: %d bytes\n", len(orderData))
	fmt.Printf("Total size: %d bytes\n", len(userData)+len(productData)+len(orderData))

	fmt.Println("✓ Size comparison completed")
	return nil
}

// createSampleOrder creates a sample order for testing (separate from manager)
func (e *Examples) createSampleOrder() *order.Order {
	now := timestamppb.Now()
	deliveryTime := timestamppb.New(time.Now().Add(5 * 24 * time.Hour))

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
				AmountCents: 7200,
			},
			ShippingCost: &product.Price{
				Currency:    "USD",
				AmountCents: 999,
			},
			Discount: &product.Price{
				Currency:    "USD",
				AmountCents: 0,
			},
			Total: &product.Price{
				Currency:    "USD",
				AmountCents: 88197,
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