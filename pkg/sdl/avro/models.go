package avro

import (
	"time"
)

// UserStatus represents the user status enum
type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusDeleted   UserStatus = "DELETED"
)

// ProductStatus represents the product status enum
type ProductStatus string

const (
	ProductStatusActive        ProductStatus = "ACTIVE"
	ProductStatusInactive      ProductStatus = "INACTIVE"
	ProductStatusOutOfStock    ProductStatus = "OUT_OF_STOCK"
	ProductStatusDiscontinued  ProductStatus = "DISCONTINUED"
)

// OrderStatus represents the order status enum
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRefunded   OrderStatus = "REFUNDED"
)

// PaymentStatus represents the payment status enum
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusAuthorized PaymentStatus = "AUTHORIZED"
	PaymentStatusCaptured   PaymentStatus = "CAPTURED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

// User represents a user entity
type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Status    UserStatus `json:"status"`
	Profile   *Profile   `json:"profile"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// Profile contains user profile information
type Profile struct {
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Phone     *string           `json:"phone"`
	Address   *Address          `json:"address"`
	Interests []string          `json:"interests"`
	Metadata  map[string]string `json:"metadata"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

// Product represents a product entity
type Product struct {
	ID            int64                 `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	SKU           string                `json:"sku"`
	Price         Price                 `json:"price"`
	Inventory     Inventory             `json:"inventory"`
	Categories    []string              `json:"categories"`
	Tags          []string              `json:"tags"`
	Status        ProductStatus         `json:"status"`
	Specifications map[string]string    `json:"specifications"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
}

// Price contains pricing information
type Price struct {
	Currency           string   `json:"currency"`
	AmountCents        int64    `json:"amountCents"`
	DiscountPercentage *float32 `json:"discountPercentage"`
}

// Inventory tracks product availability
type Inventory struct {
	Quantity       int32 `json:"quantity"`
	Reserved       int32 `json:"reserved"`
	Available      int32 `json:"available"`
	TrackInventory bool  `json:"trackInventory"`
	ReorderLevel   int32 `json:"reorderLevel"`
	MaxStock       int32 `json:"maxStock"`
}

// Order represents an order entity
type Order struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"userId"`
	OrderNumber  string        `json:"orderNumber"`
	Status       OrderStatus   `json:"status"`
	Items        []OrderItem   `json:"items"`
	Summary      OrderSummary  `json:"summary"`
	ShippingInfo *ShippingInfo `json:"shippingInfo"`
	PaymentInfo  *PaymentInfo  `json:"paymentInfo"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	ShippedAt    *time.Time    `json:"shippedAt"`
	DeliveredAt  *time.Time    `json:"deliveredAt"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID      int64             `json:"productId"`
	ProductName    string            `json:"productName"`
	ProductSKU     string            `json:"productSku"`
	Quantity       int32             `json:"quantity"`
	UnitPrice      Price             `json:"unitPrice"`
	TotalPrice     Price             `json:"totalPrice"`
	ProductVariant map[string]string `json:"productVariant"`
}

// OrderSummary contains order totals
type OrderSummary struct {
	Subtotal     Price `json:"subtotal"`
	Tax          Price `json:"tax"`
	ShippingCost Price `json:"shippingCost"`
	Discount     Price `json:"discount"`
	Total        Price `json:"total"`
	TotalItems   int32 `json:"totalItems"`
}

// ShippingInfo contains shipping details
type ShippingInfo struct {
	Address           ShippingAddress `json:"address"`
	Method            string          `json:"method"`
	TrackingNumber    *string         `json:"trackingNumber"`
	Carrier           *string         `json:"carrier"`
	Cost              Price           `json:"cost"`
	EstimatedDelivery *time.Time      `json:"estimatedDelivery"`
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	RecipientName string `json:"recipientName"`
	Street        string `json:"street"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postalCode"`
	Country       string `json:"country"`
}

// PaymentInfo contains payment details
type PaymentInfo struct {
	Method        string         `json:"method"`
	Status        PaymentStatus  `json:"status"`
	TransactionID *string        `json:"transactionId"`
	Amount        Price          `json:"amount"`
	ProcessedAt   *time.Time     `json:"processedAt"`
}

// Analytics represents analytics data
type Analytics struct {
	ID        int64             `json:"id"`
	EventType string            `json:"eventType"`
	UserID    *int64            `json:"userId"`
	SessionID string            `json:"sessionId"`
	Timestamp time.Time         `json:"timestamp"`
	Properties map[string]string `json:"properties"`
	Metrics   map[string]float64 `json:"metrics"`
	DeviceInfo *DeviceInfo       `json:"deviceInfo"`
	Location  *Location         `json:"location"`
}

// DeviceInfo contains device information
type DeviceInfo struct {
	UserAgent string `json:"userAgent"`
	Platform  string `json:"platform"`
	Browser   string `json:"browser"`
	Version   string `json:"version"`
	Mobile    bool   `json:"mobile"`
}

// Location contains geographical information
type Location struct {
	Country   string   `json:"country"`
	Region    *string  `json:"region"`
	City      *string  `json:"city"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}