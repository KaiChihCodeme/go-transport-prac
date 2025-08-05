package parquet

import (
	"time"
)

// User represents a user entity for Parquet storage
type User struct {
	ID        int64     `parquet:"id"`
	Email     string    `parquet:"email"`
	Name      string    `parquet:"name"`
	Status    string    `parquet:"status"`
	Profile   *Profile  `parquet:"profile"`
	CreatedAt time.Time `parquet:"created_at"`
	UpdatedAt time.Time `parquet:"updated_at"`
}

// Profile contains user profile information
type Profile struct {
	FirstName string            `parquet:"first_name"`
	LastName  string            `parquet:"last_name"`
	Phone     string            `parquet:"phone,optional"`
	Address   *Address          `parquet:"address,optional"`
	Interests []string          `parquet:"interests"`
	Metadata  map[string]string `parquet:"metadata"`
}

// Address represents a physical address
type Address struct {
	Street     string `parquet:"street"`
	City       string `parquet:"city"`
	State      string `parquet:"state"`
	PostalCode string `parquet:"postal_code"`
	Country    string `parquet:"country"`
}

// Product represents a product entity for Parquet storage
type Product struct {
	ID            int64                 `parquet:"id"`
	Name          string                `parquet:"name"`
	Description   string                `parquet:"description"`
	SKU           string                `parquet:"sku"`
	Price         *Price                `parquet:"price"`
	Inventory     *Inventory            `parquet:"inventory"`
	Categories    []string              `parquet:"categories"`
	Tags          []string              `parquet:"tags"`
	Status        string                `parquet:"status"`
	Specifications map[string]string    `parquet:"specifications"`
	CreatedAt     time.Time             `parquet:"created_at"`
	UpdatedAt     time.Time             `parquet:"updated_at"`
}

// Price contains pricing information
type Price struct {
	Currency           string  `parquet:"currency"`
	AmountCents        int64   `parquet:"amount_cents"`
	DiscountPercentage float32 `parquet:"discount_percentage,optional"`
}

// Inventory tracks product availability
type Inventory struct {
	Quantity       int32 `parquet:"quantity"`
	Reserved       int32 `parquet:"reserved"`
	Available      int32 `parquet:"available"`
	TrackInventory bool  `parquet:"track_inventory"`
	ReorderLevel   int32 `parquet:"reorder_level"`
	MaxStock       int32 `parquet:"max_stock"`
}

// Order represents an order entity for Parquet storage  
type Order struct {
	ID          int64        `parquet:"id,int64"`
	UserID      int64        `parquet:"user_id,int64"`
	OrderNumber string       `parquet:"order_number,utf8"`
	Status      string       `parquet:"status,utf8"`
	Items       []*OrderItem `parquet:"items,list"`
	Summary     *OrderSummary `parquet:"summary,group"`
	CreatedAt   time.Time    `parquet:"created_at,timestamp(millisecond)"`
	UpdatedAt   time.Time    `parquet:"updated_at,timestamp(millisecond)"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID   int64             `parquet:"product_id,int64"`
	ProductName string            `parquet:"product_name,utf8"`
	ProductSKU  string            `parquet:"product_sku,utf8"`
	Quantity    int32             `parquet:"quantity,int32"`
	UnitPrice   *Price            `parquet:"unit_price,group"`
	TotalPrice  *Price            `parquet:"total_price,group"`
	Variant     map[string]string `parquet:"variant,map"`
}

// OrderSummary contains order totals
type OrderSummary struct {
	Subtotal     *Price `parquet:"subtotal,group"`
	Tax          *Price `parquet:"tax,group"`
	ShippingCost *Price `parquet:"shipping_cost,group"`
	Discount     *Price `parquet:"discount,group"`
	Total        *Price `parquet:"total,group"`
	TotalItems   int32  `parquet:"total_items,int32"`
}

// Analytics represents analytics data for demonstration
type Analytics struct {
	ID            int64             `parquet:"id,int64"`
	EventType     string            `parquet:"event_type,utf8"`
	UserID        int64             `parquet:"user_id,int64,optional"`
	SessionID     string            `parquet:"session_id,utf8"`
	Timestamp     time.Time         `parquet:"timestamp,timestamp(millisecond)"`
	Properties    map[string]string `parquet:"properties,map"`
	Metrics       map[string]float64 `parquet:"metrics,map"`
	DeviceInfo    *DeviceInfo       `parquet:"device_info,group,optional"`
	Location      *Location         `parquet:"location,group,optional"`
}

// DeviceInfo contains device information
type DeviceInfo struct {
	UserAgent string `parquet:"user_agent,utf8"`
	Platform  string `parquet:"platform,utf8"`
	Browser   string `parquet:"browser,utf8,optional"`
	Version   string `parquet:"version,utf8,optional"`
	Mobile    bool   `parquet:"mobile,boolean"`
}

// Location contains geographical information
type Location struct {
	Country   string  `parquet:"country,utf8"`
	Region    string  `parquet:"region,utf8,optional"`
	City      string  `parquet:"city,utf8,optional"`
	Latitude  float64 `parquet:"latitude,double,optional"`
	Longitude float64 `parquet:"longitude,double,optional"`
}

// TimeSeriesData represents time series data for analytics
type TimeSeriesData struct {
	Timestamp time.Time `parquet:"timestamp,timestamp(millisecond)"`
	MetricName string   `parquet:"metric_name,utf8"`
	Value     float64   `parquet:"value,double"`
	Tags      map[string]string `parquet:"tags,map"`
	UserID    int64     `parquet:"user_id,int64,optional"`
	SessionID string    `parquet:"session_id,utf8,optional"`
}