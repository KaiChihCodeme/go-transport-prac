package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"go-transport-prac/internal/logger"
)

// TestHelper provides common test utilities
type TestHelper struct {
	t      *testing.T
	logger *logger.Logger
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	// Create a test logger that doesn't output to console during tests
	zapLogger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	testLogger := &logger.Logger{
		Logger: zapLogger,
	}

	return &TestHelper{
		t:      t,
		logger: testLogger,
	}
}

// Logger returns the test logger
func (h *TestHelper) Logger() *logger.Logger {
	return h.logger
}

// AssertError asserts that an error occurred and optionally checks the error message
func (h *TestHelper) AssertError(err error, expectedMessage ...string) {
	require.Error(h.t, err)
	if len(expectedMessage) > 0 {
		assert.Contains(h.t, err.Error(), expectedMessage[0])
	}
}

// AssertNoError asserts that no error occurred
func (h *TestHelper) AssertNoError(err error) {
	require.NoError(h.t, err)
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual interface{}) {
	assert.Equal(h.t, expected, actual)
}

// AssertNotEqual asserts that two values are not equal
func (h *TestHelper) AssertNotEqual(expected, actual interface{}) {
	assert.NotEqual(h.t, expected, actual)
}

// AssertTrue asserts that a condition is true
func (h *TestHelper) AssertTrue(condition bool) {
	assert.True(h.t, condition)
}

// AssertFalse asserts that a condition is false
func (h *TestHelper) AssertFalse(condition bool) {
	assert.False(h.t, condition)
}

// AssertContains asserts that a string contains a substring
func (h *TestHelper) AssertContains(str, substr string) {
	assert.Contains(h.t, str, substr)
}

// AssertJSONEqual asserts that two JSON strings are equal
func (h *TestHelper) AssertJSONEqual(expected, actual string) {
	var expectedMap, actualMap interface{}
	require.NoError(h.t, json.Unmarshal([]byte(expected), &expectedMap))
	require.NoError(h.t, json.Unmarshal([]byte(actual), &actualMap))
	assert.Equal(h.t, expectedMap, actualMap)
}

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	*TestHelper
	server *httptest.Server
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper(t *testing.T, handler http.Handler) *HTTPTestHelper {
	return &HTTPTestHelper{
		TestHelper: NewTestHelper(t),
		server:     httptest.NewServer(handler),
	}
}

// Close closes the test server
func (h *HTTPTestHelper) Close() {
	h.server.Close()
}

// URL returns the test server URL
func (h *HTTPTestHelper) URL() string {
	return h.server.URL
}

// GET performs a GET request
func (h *HTTPTestHelper) GET(path string, headers ...map[string]string) *http.Response {
	return h.request("GET", path, nil, headers...)
}

// POST performs a POST request
func (h *HTTPTestHelper) POST(path string, body interface{}, headers ...map[string]string) *http.Response {
	return h.request("POST", path, body, headers...)
}

// PUT performs a PUT request
func (h *HTTPTestHelper) PUT(path string, body interface{}, headers ...map[string]string) *http.Response {
	return h.request("PUT", path, body, headers...)
}

// DELETE performs a DELETE request
func (h *HTTPTestHelper) DELETE(path string, headers ...map[string]string) *http.Response {
	return h.request("DELETE", path, nil, headers...)
}

// PATCH performs a PATCH request
func (h *HTTPTestHelper) PATCH(path string, body interface{}, headers ...map[string]string) *http.Response {
	return h.request("PATCH", path, body, headers...)
}

// request performs an HTTP request
func (h *HTTPTestHelper) request(method, path string, body interface{}, headers ...map[string]string) *http.Response {
	var bodyReader io.Reader

	if body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		case []byte:
			bodyReader = bytes.NewReader(v)
		default:
			jsonData, err := json.Marshal(body)
			require.NoError(h.t, err)
			bodyReader = bytes.NewReader(jsonData)
		}
	}

	url := h.server.URL + path
	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(h.t, err)

	// Set default content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set additional headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(h.t, err)

	return resp
}

// AssertStatusCode asserts the HTTP status code
func (h *HTTPTestHelper) AssertStatusCode(resp *http.Response, expectedCode int) {
	assert.Equal(h.t, expectedCode, resp.StatusCode)
}

// AssertResponseBody asserts the response body content
func (h *HTTPTestHelper) AssertResponseBody(resp *http.Response, expected string) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(h.t, err)
	defer resp.Body.Close()
	assert.Equal(h.t, expected, string(body))
}

// AssertResponseJSON asserts the response body as JSON
func (h *HTTPTestHelper) AssertResponseJSON(resp *http.Response, expected interface{}) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(h.t, err)
	defer resp.Body.Close()

	expectedJSON, err := json.Marshal(expected)
	require.NoError(h.t, err)

	h.AssertJSONEqual(string(expectedJSON), string(body))
}

// GetResponseBody returns the response body as string
func (h *HTTPTestHelper) GetResponseBody(resp *http.Response) string {
	body, err := io.ReadAll(resp.Body)
	require.NoError(h.t, err)
	defer resp.Body.Close()
	return string(body)
}

// FileTestHelper provides utilities for file-based testing
type FileTestHelper struct {
	*TestHelper
	tempDir string
}

// NewFileTestHelper creates a new file test helper
func NewFileTestHelper(t *testing.T) *FileTestHelper {
	tempDir, err := os.MkdirTemp("", "transport-test-*")
	require.NoError(t, err)

	helper := &FileTestHelper{
		TestHelper: NewTestHelper(t),
		tempDir:    tempDir,
	}

	// Cleanup temp directory when test finishes
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return helper
}

// TempDir returns the temporary directory path
func (h *FileTestHelper) TempDir() string {
	return h.tempDir
}

// CreateTempFile creates a temporary file with content
func (h *FileTestHelper) CreateTempFile(name, content string) string {
	filePath := filepath.Join(h.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(h.t, err)
	return filePath
}

// CreateTempDir creates a temporary directory
func (h *FileTestHelper) CreateTempDir(name string) string {
	dirPath := filepath.Join(h.tempDir, name)
	err := os.MkdirAll(dirPath, 0755)
	require.NoError(h.t, err)
	return dirPath
}

// ReadFile reads a file and returns its content
func (h *FileTestHelper) ReadFile(path string) string {
	content, err := os.ReadFile(path)
	require.NoError(h.t, err)
	return string(content)
}

// AssertFileExists asserts that a file exists
func (h *FileTestHelper) AssertFileExists(path string) {
	_, err := os.Stat(path)
	assert.NoError(h.t, err, "File should exist: %s", path)
}

// AssertFileNotExists asserts that a file does not exist
func (h *FileTestHelper) AssertFileNotExists(path string) {
	_, err := os.Stat(path)
	assert.True(h.t, os.IsNotExist(err), "File should not exist: %s", path)
}

// AssertFileContent asserts the content of a file
func (h *FileTestHelper) AssertFileContent(path, expectedContent string) {
	actualContent := h.ReadFile(path)
	assert.Equal(h.t, expectedContent, actualContent)
}

// TimeoutContext creates a context with timeout for testing
func TimeoutContext(t *testing.T, timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}

// WaitForCondition waits for a condition to become true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Condition not met within %v: %s", timeout, message)
}

// MockData provides common test data
var MockData = struct {
	UserID        string
	RequestID     string
	ComponentName string
	JSONData      string
	XMLData       string
	CSVData       string
}{
	UserID:        "test-user-123",
	RequestID:     "req-456",
	ComponentName: "test-component",
	JSONData: `{
		"id": 1,
		"name": "Test User",
		"email": "test@example.com",
		"active": true,
		"metadata": {
			"created_at": "2023-01-01T00:00:00Z",
			"tags": ["test", "user"]
		}
	}`,
	XMLData: `<?xml version="1.0" encoding="UTF-8"?>
	<user>
		<id>1</id>
		<name>Test User</name>
		<email>test@example.com</email>
		<active>true</active>
	</user>`,
	CSVData: `id,name,email,active
1,Test User,test@example.com,true
2,Another User,another@example.com,false`,
}

// TableTest represents a table-driven test case
type TableTest struct {
	Name     string
	Input    interface{}
	Expected interface{}
	Error    string
	Setup    func()
	Cleanup  func()
}

// RunTableTests runs table-driven tests
func RunTableTests(t *testing.T, tests []TableTest, testFunc func(*testing.T, TableTest)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Setup != nil {
				tt.Setup()
			}
			if tt.Cleanup != nil {
				defer tt.Cleanup()
			}
			testFunc(t, tt)
		})
	}
}

// BenchmarkHelper provides utilities for benchmarking
type BenchmarkHelper struct {
	b *testing.B
}

// NewBenchmarkHelper creates a new benchmark helper
func NewBenchmarkHelper(b *testing.B) *BenchmarkHelper {
	return &BenchmarkHelper{b: b}
}

// ResetTimer resets the benchmark timer
func (h *BenchmarkHelper) ResetTimer() {
	h.b.ResetTimer()
}

// StartTimer starts the benchmark timer
func (h *BenchmarkHelper) StartTimer() {
	h.b.StartTimer()
}

// StopTimer stops the benchmark timer
func (h *BenchmarkHelper) StopTimer() {
	h.b.StopTimer()
}

// ReportAllocs enables allocation reporting
func (h *BenchmarkHelper) ReportAllocs() {
	h.b.ReportAllocs()
}

// SetBytes sets the number of bytes processed per iteration
func (h *BenchmarkHelper) SetBytes(n int64) {
	h.b.SetBytes(n)
}

// Logf logs a message during benchmarking
func (h *BenchmarkHelper) Logf(format string, args ...interface{}) {
	h.b.Logf(format, args...)
}

// RunParallel runs the benchmark in parallel
func (h *BenchmarkHelper) RunParallel(body func(*testing.PB)) {
	h.b.RunParallel(body)
}

// SkipLongTest skips tests that take a long time unless explicitly enabled
func SkipLongTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long test in short mode")
	}
}

// RequireEnv requires an environment variable to be set
func RequireEnv(t *testing.T, key string) string {
	value := os.Getenv(key)
	if value == "" {
		t.Skipf("Environment variable %s is required but not set", key)
	}
	return value
}

// SetEnv sets an environment variable for the duration of the test
func SetEnv(t *testing.T, key, value string) {
	oldValue := os.Getenv(key)
	err := os.Setenv(key, value)
	require.NoError(t, err)

	t.Cleanup(func() {
		if oldValue != "" {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	})
}