package types

import (
	"context"
	"io"
	"time"
)

// Service represents a generic service interface
type Service interface {
	// Initialize initializes the service with configuration
	Initialize(ctx context.Context, config any) error
	
	// Start starts the service
	Start(ctx context.Context) error
	
	// Stop stops the service gracefully
	Stop(ctx context.Context) error
	
	// Health returns the health status of the service
	Health(ctx context.Context) error
	
	// Name returns the service name
	Name() string
}

// Repository represents a generic repository interface
type Repository interface {
	// Connect establishes connection to the data store
	Connect(ctx context.Context) error
	
	// Disconnect closes connection to the data store
	Disconnect(ctx context.Context) error
	
	// Ping checks if the connection is alive
	Ping(ctx context.Context) error
	
	// Migrate runs any necessary migrations
	Migrate(ctx context.Context) error
}

// Serializer represents a data serialization interface
type Serializer interface {
	// Serialize converts data to bytes
	Serialize(data any) ([]byte, error)
	
	// Deserialize converts bytes to data
	Deserialize(data []byte, target any) error
	
	// ContentType returns the content type of the serialized data
	ContentType() string
	
	// FileExtension returns the file extension for this format
	FileExtension() string
}

// SchemaValidator represents a schema validation interface
type SchemaValidator interface {
	// Validate validates data against a schema
	Validate(data any) error
	
	// ValidateBytes validates byte data against a schema
	ValidateBytes(data []byte) error
	
	// SetSchema sets the validation schema
	SetSchema(schema any) error
	
	// GetSchema returns the current schema
	GetSchema() any
}

// Encoder represents a data encoding interface
type Encoder interface {
	// Encode encodes data to a writer
	Encode(w io.Writer, data any) error
	
	// EncodeToBytes encodes data to bytes
	EncodeToBytes(data any) ([]byte, error)
}

// Decoder represents a data decoding interface
type Decoder interface {
	// Decode decodes data from a reader
	Decode(r io.Reader, target any) error
	
	// DecodeFromBytes decodes data from bytes
	DecodeFromBytes(data []byte, target any) error
}

// Codec combines encoder and decoder interfaces
type Codec interface {
	Encoder
	Decoder
}

// HTTPHandler represents an HTTP request handler
type HTTPHandler interface {
	// Handle processes an HTTP request
	Handle(ctx context.Context, request HTTPRequest) (HTTPResponse, error)
	
	// Method returns the HTTP method this handler supports
	Method() string
	
	// Path returns the URL path this handler supports
	Path() string
}

// HTTPRequest represents an HTTP request
type HTTPRequest struct {
	Method     string
	Path       string
	Headers    map[string]string
	Body       []byte
	Query      map[string]string
	PathParams map[string]string
	UserID     string
	RequestID  string
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// WebSocketHandler represents a WebSocket message handler
type WebSocketHandler interface {
	// OnConnect handles new WebSocket connections
	OnConnect(ctx context.Context, conn WebSocketConnection) error
	
	// OnMessage handles incoming messages
	OnMessage(ctx context.Context, conn WebSocketConnection, message []byte) error
	
	// OnDisconnect handles connection disconnections
	OnDisconnect(ctx context.Context, conn WebSocketConnection) error
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection interface {
	// Send sends a message to the connection
	Send(ctx context.Context, message []byte) error
	
	// Close closes the connection
	Close() error
	
	// ID returns the connection ID
	ID() string
	
	// UserID returns the user ID associated with this connection
	UserID() string
}

// MessageBroker represents a message broker interface
type MessageBroker interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message []byte) error
	
	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	
	// Unsubscribe unsubscribes from a topic
	Unsubscribe(ctx context.Context, topic string) error
	
	// Close closes the broker connection
	Close() error
}

// MessageHandler represents a message handler function
type MessageHandler func(ctx context.Context, message Message) error

// Message represents a message
type Message struct {
	ID        string
	Topic     string
	Data      []byte
	Headers   map[string]string
	Timestamp time.Time
}

// Cache represents a caching interface
type Cache interface {
	// Get retrieves a value by key
	Get(ctx context.Context, key string) ([]byte, error)
	
	// Set stores a value with key and expiration
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	
	// Delete removes a value by key
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)
	
	// Close closes the cache connection
	Close() error
}

// Storage represents a file/object storage interface
type Storage interface {
	// Put stores data with the given key
	Put(ctx context.Context, key string, data io.Reader) error
	
	// Get retrieves data by key
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	
	// Delete removes data by key
	Delete(ctx context.Context, key string) error
	
	// Exists checks if data exists for the given key
	Exists(ctx context.Context, key string) (bool, error)
	
	// List lists all keys with the given prefix
	List(ctx context.Context, prefix string) ([]string, error)
}

// MetricsCollector represents a metrics collection interface
type MetricsCollector interface {
	// Counter increments a counter metric
	Counter(name string, tags map[string]string, value float64)
	
	// Gauge sets a gauge metric
	Gauge(name string, tags map[string]string, value float64)
	
	// Histogram records a histogram metric
	Histogram(name string, tags map[string]string, value float64)
	
	// Timer records a timing metric
	Timer(name string, tags map[string]string, duration time.Duration)
}

// HealthChecker represents a health check interface
type HealthChecker interface {
	// Check performs a health check
	Check(ctx context.Context) error
	
	// Name returns the name of the health check
	Name() string
}

// Migrator represents a database migration interface
type Migrator interface {
	// Up runs pending migrations
	Up(ctx context.Context) error
	
	// Down rolls back the last migration
	Down(ctx context.Context) error
	
	// Version returns the current migration version
	Version(ctx context.Context) (int, error)
	
	// SetVersion sets the migration version
	SetVersion(ctx context.Context, version int) error
}

// EventEmitter represents an event emission interface
type EventEmitter interface {
	// Emit emits an event
	Emit(ctx context.Context, event Event) error
	
	// Subscribe subscribes to events
	Subscribe(ctx context.Context, eventType string, handler EventHandler) error
	
	// Unsubscribe unsubscribes from events
	Unsubscribe(ctx context.Context, eventType string, handler EventHandler) error
}

// EventHandler represents an event handler function
type EventHandler func(ctx context.Context, event Event) error

// Event represents an event
type Event struct {
	ID        string
	Type      string
	Source    string
	Data      any
	Timestamp time.Time
	Metadata  map[string]any
}

// Validator represents a generic validation interface
type Validator interface {
	// Validate validates the given data
	Validate(data any) error
	
	// ValidateStruct validates a struct
	ValidateStruct(s any) error
}

// Logger represents a logging interface
type Logger interface {
	// Debug logs a debug message
	Debug(msg string, fields ...any)
	
	// Info logs an info message
	Info(msg string, fields ...any)
	
	// Warn logs a warning message
	Warn(msg string, fields ...any)
	
	// Error logs an error message
	Error(msg string, fields ...any)
	
	// Fatal logs a fatal message and exits
	Fatal(msg string, fields ...any)
	
	// WithFields returns a logger with additional fields
	WithFields(fields map[string]any) Logger
}

// Configurable represents something that can be configured
type Configurable interface {
	// Configure configures the component with the given config
	Configure(config any) error
	
	// GetConfig returns the current configuration
	GetConfig() any
}

// Startable represents something that can be started
type Startable interface {
	// Start starts the component
	Start(ctx context.Context) error
}

// Stoppable represents something that can be stopped
type Stoppable interface {
	// Stop stops the component
	Stop(ctx context.Context) error
}

// Lifecycle combines Startable and Stoppable
type Lifecycle interface {
	Startable
	Stoppable
}

// Named represents something with a name
type Named interface {
	// Name returns the name
	Name() string
}

// Versioned represents something with a version
type Versioned interface {
	// Version returns the version
	Version() string
}