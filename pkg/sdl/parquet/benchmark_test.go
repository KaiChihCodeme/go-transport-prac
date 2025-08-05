package parquet

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

// UserJSON for JSON comparison
type UserJSON struct {
	ID        int64       `json:"id"`
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	Profile   ProfileJSON `json:"profile"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
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

// createSampleUsers creates sample users for benchmarking
func createSampleUsers(count int) []User {
	users := make([]User, count)
	now := time.Now()
	
	for i := 0; i < count; i++ {
		users[i] = User{
			ID:     int64(i + 1),
			Email:  "benchmark@example.com",
			Name:   "Benchmark User",
			Status: "active",
			Profile: &Profile{
				FirstName: "Benchmark",
				LastName:  "User",
				Phone:     "+1-555-BENCH",
				Address: &Address{
					Street:     "123 Benchmark St",
					City:       "Test City",
					State:      "TS",
					PostalCode: "12345",
					Country:    "USA",
				},
				Interests: []string{"performance", "testing"},
				Metadata: map[string]string{
					"source": "benchmark",
					"type":   "test",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	
	return users
}

// createSampleUsersJSON creates JSON equivalent for comparison
func createSampleUsersJSON(count int) []UserJSON {
	users := make([]UserJSON, count)
	now := time.Now()
	
	for i := 0; i < count; i++ {
		users[i] = UserJSON{
			ID:     int64(i + 1),
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
					State:      "TS",
					PostalCode: "12345",
					Country:    "USA",
				},
				Interests: []string{"performance", "testing"},
				Metadata: map[string]string{
					"source": "benchmark",
					"type":   "test",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	
	return users
}

// Serialization benchmarks
func BenchmarkParquetUserSerialization(b *testing.B) {
	testDir := "tmp/bench_parquet"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(1000)
	filename := "bench_users.parquet"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		err := manager.WriteUsers(filename, users)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUserSerialization(b *testing.B) {
	users := createSampleUsersJSON(1000)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(users)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Deserialization benchmarks
func BenchmarkParquetUserDeserialization(b *testing.B) {
	testDir := "tmp/bench_parquet_read"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(1000)
	filename := "bench_read_users.parquet"
	
	// Pre-create the file
	err := manager.WriteUsers(filename, users)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, err := manager.ReadUsers(filename)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUserDeserialization(b *testing.B) {
	users := createSampleUsersJSON(1000)
	data, err := json.Marshal(users)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		var result []UserJSON
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Size comparison benchmark
func BenchmarkParquetVsJSONSize(b *testing.B) {
	testDir := "tmp/bench_size"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(1000)
	usersJSON := createSampleUsersJSON(1000)
	
	// Get Parquet size
	filename := "size_test.parquet"
	err := manager.WriteUsers(filename, users)
	if err != nil {
		b.Fatal(err)
	}
	
	info, err := manager.GetBasicFileInfo(filename)
	if err != nil {
		b.Fatal(err)
	}
	
	// Get JSON size
	jsonData, err := json.Marshal(usersJSON)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ReportMetric(float64(info.FileSize), "parquet-bytes")
	b.ReportMetric(float64(len(jsonData)), "json-bytes")
	b.ReportMetric(float64(len(jsonData))/float64(info.FileSize), "size-ratio")
}

// Full cycle benchmarks
func BenchmarkParquetFullCycle(b *testing.B) {
	testDir := "tmp/bench_full_parquet"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(100)
	filename := "full_cycle.parquet"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Write
		err := manager.WriteUsers(filename, users)  
		if err != nil {
			b.Fatal(err)
		}
		
		// Read
		_, err = manager.ReadUsers(filename)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONFullCycle(b *testing.B) {
	users := createSampleUsersJSON(100)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		// Marshal
		data, err := json.Marshal(users)
		if err != nil {
			b.Fatal(err)
		}
		
		// Unmarshal
		var result []UserJSON
		err = json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Different data sizes
func BenchmarkParquetSmallDataset(b *testing.B) {
	testDir := "tmp/bench_small"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(10)
	filename := "small.parquet"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := manager.WriteUsers(filename, users)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportMetric(10, "records")
}

func BenchmarkParquetMediumDataset(b *testing.B) {
	testDir := "tmp/bench_medium"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(1000)
	filename := "medium.parquet"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := manager.WriteUsers(filename, users)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportMetric(1000, "records")
}

func BenchmarkParquetLargeDataset(b *testing.B) {
	testDir := "tmp/bench_large"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(10000)
	filename := "large.parquet"
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		err := manager.WriteUsers(filename, users)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportMetric(10000, "records")
}

// Memory usage benchmark
func BenchmarkParquetMemoryUsage(b *testing.B) {
	testDir := "tmp/bench_memory"
	manager := NewSimpleManager(testDir)
	defer os.RemoveAll(testDir)

	users := createSampleUsers(5000)
	filename := "memory_test.parquet"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		err := manager.WriteUsers(filename, users)
		if err != nil {
			b.Fatal(err)
		}
		
		_, err = manager.ReadUsers(filename)
		if err != nil {
			b.Fatal(err)
		}
	}
}