package protobuf

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-transport-prac/pkg/sdl/protobuf/gen/order"
	"go-transport-prac/pkg/sdl/protobuf/gen/product"
	"go-transport-prac/pkg/sdl/protobuf/gen/user"
)

// createSampleUser creates a sample user for benchmarking
func createSampleUser() *user.User {
	return &user.User{
		Id:     12345,
		Email:  "benchmark@example.com",
		Name:   "Benchmark User",
		Status: user.UserStatus_USER_STATUS_ACTIVE,
		Profile: &user.Profile{
			FirstName: "Benchmark",
			LastName:  "User",
			Phone:     "+1-555-BENCH",
			Address: &user.Address{
				Street:     "123 Benchmark St",
				City:       "Test City",
				State:      "TC",
				PostalCode: "12345",
				Country:    "USA",
			},
			Interests: []string{"performance", "testing", "optimization"},
			Metadata: map[string]string{
				"preferred_language": "en",
				"timezone":           "UTC",
				"theme":              "dark",
			},
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
}

// createSampleProduct creates a sample product for benchmarking
func createSampleProduct() *product.Product {
	return &product.Product{
		Id:          67890,
		Name:        "Benchmark Product",
		Description: "A product designed for performance testing and benchmarking purposes",
		Sku:         "BENCH-001",
		Categories:  []string{"Testing", "Performance"},
		Status:      product.ProductStatus_PRODUCT_STATUS_ACTIVE,
		Price: &product.Price{
			Currency:    "USD",
			AmountCents: 9999, // $99.99 in cents
		},
		Inventory: &product.Inventory{
			Quantity:       1000,
			Reserved:       50,
			Available:      950,
			TrackInventory: true,
			ReorderLevel:   100,
			MaxStock:       2000,
		},
		Specifications: &product.Specifications{
			Attributes: map[string]string{
				"color":    "blue",
				"size":     "medium",
				"material": "plastic",
			},
			Dimensions: &product.Dimensions{
				Length: 10.5,
				Width:  5.2,
				Height: 3.1,
				Unit:   "cm",
			},
			Weight: &product.Weight{
				Value: 1.5,
				Unit:  "kg",
			},
		},
		Tags:      []string{"benchmark", "testing", "performance"},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
}

// createSampleOrder creates a sample order for benchmarking
func createSampleOrder() *order.Order {
	return &order.Order{
		Id:          11111,
		UserId:      12345,
		OrderNumber: "ORD-BENCH-001",
		Status:      order.OrderStatus_ORDER_STATUS_CONFIRMED,
		Items: []*order.OrderItem{
			{
				ProductId:   67890,
				ProductName: "Benchmark Product",
				ProductSku:  "BENCH-001",
				Quantity:    2,
				UnitPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 9999,
				},
				TotalPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 19998,
				},
			},
			{
				ProductId:   67891,
				ProductName: "Test Product 2",
				ProductSku:  "BENCH-002",
				Quantity:    1,
				UnitPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 4999,
				},
				TotalPrice: &product.Price{
					Currency:    "USD",
					AmountCents: 4999,
				},
			},
		},
		Summary: &order.OrderSummary{
			Subtotal: &product.Price{
				Currency:    "USD",
				AmountCents: 24997,
			},
			Tax: &product.Price{
				Currency:    "USD",
				AmountCents: 2000,
			},
			ShippingCost: &product.Price{
				Currency:    "USD",
				AmountCents: 500,
			},
			Discount: &product.Price{
				Currency:    "USD",
				AmountCents: 0,
			},
			Total: &product.Price{
				Currency:    "USD",
				AmountCents: 27497,
			},
			TotalItems: 3,
		},
		Shipping: &order.ShippingInfo{
			Address: &user.Address{
				Street:     "123 Benchmark St",
				City:       "Test City",
				State:      "TC",
				PostalCode: "12345",
				Country:    "USA",
			},
			Method:         "standard",
			TrackingNumber: "TRACK123456",
			Carrier:        "UPS",
			Cost: &product.Price{
				Currency:    "USD",
				AmountCents: 500,
			},
			EstimatedDelivery: timestamppb.Now(),
		},
		Payment: &order.PaymentInfo{
			Method:        "credit_card",
			Status:        order.PaymentStatus_PAYMENT_STATUS_CAPTURED,
			TransactionId: "TXN123456789",
			Amount: &product.Price{
				Currency:    "USD",
				AmountCents: 27497,
			},
			ProcessedAt: timestamppb.Now(),
		},
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
}

// JSON equivalents for comparison
type UserJSON struct {
	ID        uint64      `json:"id"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Profile   ProfileJSON `json:"profile"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type ProfileJSON struct {
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Phone     string            `json:"phone"`
	Address   AddressJSON       `json:"address"`
	Interests []string          `json:"interests"`
	Metadata  map[string]string `json:"metadata"`
}

type AddressJSON struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

func createSampleUserJSON() UserJSON {
	return UserJSON{
		ID:     12345,
		Email:  "benchmark@example.com",
		Name:   "Benchmark User",
		Status: "active",
		Profile: ProfileJSON{
			FirstName: "Benchmark",
			LastName:  "User",
			Phone:     "+1-555-BENCH",
			Address: AddressJSON{
				Street:     "123 Benchmark St",
				City:       "Test City",
				State:      "TC",
				PostalCode: "12345",
				Country:    "USA",
			},
			Interests: []string{"performance", "testing", "optimization"},
			Metadata: map[string]string{
				"preferred_language": "en",
				"timezone":           "UTC",
				"theme":              "dark",
			},
		},
		CreatedAt: "2023-01-01T00:00:00Z",
		UpdatedAt: "2023-01-01T00:00:00Z",
	}
}

// Benchmarks for serialization
func BenchmarkProtobufUserSerialization(b *testing.B) {
	user := createSampleUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUserSerialization(b *testing.B) {
	user := createSampleUserJSON()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtobufProductSerialization(b *testing.B) {
	product := createSampleProduct()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(product)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtobufOrderSerialization(b *testing.B) {
	order := createSampleOrder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(order)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarks for deserialization
func BenchmarkProtobufUserDeserialization(b *testing.B) {
	sampleUser := createSampleUser()
	data, err := proto.Marshal(sampleUser)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var u user.User
		err := proto.Unmarshal(data, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUserDeserialization(b *testing.B) {
	user := createSampleUserJSON()
	data, err := json.Marshal(user)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var u UserJSON
		err := json.Unmarshal(data, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtobufProductDeserialization(b *testing.B) {
	sampleProduct := createSampleProduct()
	data, err := proto.Marshal(sampleProduct)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var p product.Product
		err := proto.Unmarshal(data, &p)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtobufOrderDeserialization(b *testing.B) {
	sampleOrder := createSampleOrder()
	data, err := proto.Marshal(sampleOrder)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var o order.Order
		err := proto.Unmarshal(data, &o)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Memory and size benchmarks
func BenchmarkProtobufVsJSONSize(b *testing.B) {
	user := createSampleUser()
	userJSON := createSampleUserJSON()

	protoData, err := proto.Marshal(user)
	if err != nil {
		b.Fatal(err)
	}

	jsonData, err := json.Marshal(userJSON)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportMetric(float64(len(protoData)), "protobuf-bytes")
	b.ReportMetric(float64(len(jsonData)), "json-bytes")
	b.ReportMetric(float64(len(jsonData))/float64(len(protoData)), "size-ratio")
}

// Comprehensive benchmark comparing all operations
func BenchmarkProtobufFullCycle(b *testing.B) {
	sampleUser := createSampleUser()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Serialize
		data, err := proto.Marshal(sampleUser)
		if err != nil {
			b.Fatal(err)
		}

		// Deserialize
		var u user.User
		err = proto.Unmarshal(data, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONFullCycle(b *testing.B) {
	user := createSampleUserJSON()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Serialize
		data, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}

		// Deserialize
		var u UserJSON
		err = json.Unmarshal(data, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Large data set benchmarks
func BenchmarkProtobufLargeDataSet(b *testing.B) {
	// Create 1000 users
	users := make([]*user.User, 1000)
	for i := 0; i < 1000; i++ {
		users[i] = createSampleUser()
		users[i].Id = uint64(i + 1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, u := range users {
			data, err := proto.Marshal(u)
			if err != nil {
				b.Fatal(err)
			}

			var newUser user.User
			err = proto.Unmarshal(data, &newUser)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.ReportMetric(1000, "users-processed")
}
