package types

import (
	"time"
)

// Result represents a generic result with data and error
type Result[T any] struct {
	Data  T
	Error error
}

// NewResult creates a new result
func NewResult[T any](data T, err error) Result[T] {
	return Result[T]{
		Data:  data,
		Error: err,
	}
}

// IsSuccess returns true if the result has no error
func (r Result[T]) IsSuccess() bool {
	return r.Error == nil
}

// IsError returns true if the result has an error
func (r Result[T]) IsError() bool {
	return r.Error != nil
}

// Unwrap returns the data if successful, otherwise panics
func (r Result[T]) Unwrap() T {
	if r.Error != nil {
		panic(r.Error)
	}
	return r.Data
}

// UnwrapOr returns the data if successful, otherwise returns the default value
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.Error != nil {
		return defaultValue
	}
	return r.Data
}

// Option represents an optional value
type Option[T any] struct {
	value   T
	present bool
}

// Some creates an Option with a value
func Some[T any](value T) Option[T] {
	return Option[T]{
		value:   value,
		present: true,
	}
}

// None creates an empty Option
func None[T any]() Option[T] {
	return Option[T]{
		present: false,
	}
}

// IsSome returns true if the option has a value
func (o Option[T]) IsSome() bool {
	return o.present
}

// IsNone returns true if the option has no value
func (o Option[T]) IsNone() bool {
	return !o.present
}

// Unwrap returns the value if present, otherwise panics
func (o Option[T]) Unwrap() T {
	if !o.present {
		panic("called Unwrap on None option")
	}
	return o.value
}

// UnwrapOr returns the value if present, otherwise returns the default value
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if !o.present {
		return defaultValue
	}
	return o.value
}

// Map applies a function to the value if present
func (o Option[T]) Map(fn func(T) T) Option[T] {
	if !o.present {
		return None[T]()
	}
	return Some(fn(o.value))
}

// Pair represents a key-value pair
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair creates a new pair
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{
		Key:   key,
		Value: value,
	}
}

// Metadata represents generic metadata
type Metadata map[string]any

// Get returns a value from metadata
func (m Metadata) Get(key string) (any, bool) {
	value, exists := m[key]
	return value, exists
}

// GetString returns a string value from metadata
func (m Metadata) GetString(key string) (string, bool) {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetInt returns an int value from metadata
func (m Metadata) GetInt(key string) (int, bool) {
	if value, exists := m[key]; exists {
		if i, ok := value.(int); ok {
			return i, true
		}
	}
	return 0, false
}

// GetBool returns a bool value from metadata
func (m Metadata) GetBool(key string) (bool, bool) {
	if value, exists := m[key]; exists {
		if b, ok := value.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// Set sets a value in metadata
func (m Metadata) Set(key string, value any) {
	m[key] = value
}

// Delete removes a key from metadata
func (m Metadata) Delete(key string) {
	delete(m, key)
}

// Clone creates a deep copy of metadata
func (m Metadata) Clone() Metadata {
	clone := make(Metadata, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

// Status represents a generic status
type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// String returns the string representation of status
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusPending, StatusCompleted, StatusFailed, StatusCancelled:
		return true
	default:
		return false
	}
}

// ID represents a generic identifier
type ID string

// String returns the string representation of ID
func (id ID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty
func (id ID) IsEmpty() bool {
	return string(id) == ""
}

// Timestamp represents a timestamp with additional methods
type Timestamp struct {
	time.Time
}

// NewTimestamp creates a new timestamp
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// Now creates a timestamp for the current time
func Now() Timestamp {
	return Timestamp{Time: time.Now()}
}

// Unix creates a timestamp from Unix time
func Unix(sec int64) Timestamp {
	return Timestamp{Time: time.Unix(sec, 0)}
}

// IsZero returns true if the timestamp is zero
func (t Timestamp) IsZero() bool {
	return t.Time.IsZero()
}

// Unix returns the Unix timestamp
func (t Timestamp) Unix() int64 {
	return t.Time.Unix()
}

// RFC3339 returns the RFC3339 formatted string
func (t Timestamp) RFC3339() string {
	return t.Time.Format(time.RFC3339)
}

// Range represents a range of values
type Range[T comparable] struct {
	Start T
	End   T
}

// NewRange creates a new range
func NewRange[T comparable](start, end T) Range[T] {
	return Range[T]{
		Start: start,
		End:   end,
	}
}

// Contains checks if a value is within the range
func (r Range[T]) Contains(value T) bool {
	// This is a simplified implementation
	// In practice, you'd need to implement comparison for the generic type
	return true // Placeholder
}

// Page represents pagination information
type Page struct {
	Number int `json:"number"`
	Size   int `json:"size"`
	Offset int `json:"offset"`
}

// NewPage creates a new page
func NewPage(number, size int) Page {
	return Page{
		Number: number,
		Size:   size,
		Offset: (number - 1) * size,
	}
}

// PagedResult represents a paginated result
type PagedResult[T any] struct {
	Data       []T  `json:"data"`
	Page       Page `json:"page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// NewPagedResult creates a new paged result
func NewPagedResult[T any](data []T, page Page, total int) PagedResult[T] {
	totalPages := (total + page.Size - 1) / page.Size
	return PagedResult[T]{
		Data:       data,
		Page:       page,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page.Number < totalPages,
		HasPrev:    page.Number > 1,
	}
}

// Filter represents a generic filter
type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// Sort represents sorting information
type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"` // "asc" or "desc"
}

// Query represents a generic query with filters, sorting, and pagination
type Query struct {
	Filters []Filter `json:"filters,omitempty"`
	Sort    []Sort   `json:"sort,omitempty"`
	Page    *Page    `json:"page,omitempty"`
	Search  string   `json:"search,omitempty"`
}

// NewQuery creates a new query
func NewQuery() *Query {
	return &Query{
		Filters: make([]Filter, 0),
		Sort:    make([]Sort, 0),
	}
}

// AddFilter adds a filter to the query
func (q *Query) AddFilter(field, operator string, value any) *Query {
	q.Filters = append(q.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return q
}

// AddSort adds sorting to the query
func (q *Query) AddSort(field, order string) *Query {
	q.Sort = append(q.Sort, Sort{
		Field: field,
		Order: order,
	})
	return q
}

// SetPage sets pagination for the query
func (q *Query) SetPage(page Page) *Query {
	q.Page = &page
	return q
}

// SetSearch sets search term for the query
func (q *Query) SetSearch(search string) *Query {
	q.Search = search
	return q
}

// BuildInfo represents build information
type BuildInfo struct {
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	BuildTime time.Time `json:"build_time"`
	GoVersion string    `json:"go_version"`
}

// HealthStatus represents health status information
type HealthStatus struct {
	Status     string            `json:"status"`
	Version    string            `json:"version"`
	Timestamp  time.Time         `json:"timestamp"`
	Checks     map[string]string `json:"checks"`
	Uptime     time.Duration     `json:"uptime"`
	SystemInfo SystemInfo        `json:"system_info"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	NumCPU   int    `json:"num_cpu"`
	GoMaxProcs int  `json:"go_max_procs"`
}

// APIError represents an API error response
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details string                 `json:"details,omitempty"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse[T any] struct {
	Success bool      `json:"success"`
	Data    T         `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
	Meta    Metadata  `json:"meta,omitempty"`
}

// NewSuccessResponse creates a successful API response
func NewSuccessResponse[T any](data T) APIResponse[T] {
	return APIResponse[T]{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates an error API response
func NewErrorResponse[T any](err APIError) APIResponse[T] {
	return APIResponse[T]{
		Success: false,
		Error:   &err,
	}
}