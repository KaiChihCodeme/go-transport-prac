package avro

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// BenchmarkResults contains performance comparison results
type BenchmarkResults struct {
	Format               string        `json:"format"`
	SerializationTime    time.Duration `json:"serializationTime"`
	DeserializationTime  time.Duration `json:"deserializationTime"`
	SerializedSize       int           `json:"serializedSize"`
	MemoryUsage          int64         `json:"memoryUsage"`
	ItemsPerSecond       float64       `json:"itemsPerSecond"`
}

// PerformanceBenchmark runs performance tests comparing different serialization formats
type PerformanceBenchmark struct {
	manager *Manager
	users   []User
	products []Product
}

// NewPerformanceBenchmark creates a new performance benchmark
func NewPerformanceBenchmark() (*PerformanceBenchmark, error) {
	manager, err := NewManager("tmp/benchmark")
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	pb := &PerformanceBenchmark{
		manager: manager,
	}

	// Generate test data
	pb.users = manager.CreateSampleUsers(1000)
	pb.products = manager.CreateSampleProducts(1000)

	return pb, nil
}

// RunBenchmarks executes all performance benchmarks
func (pb *PerformanceBenchmark) RunBenchmarks() error {
	fmt.Println("=== Performance Benchmarks ===")
	fmt.Printf("Testing with %d users and %d products\n", len(pb.users), len(pb.products))

	// Run user benchmarks
	fmt.Println("--- User Serialization Benchmarks ---")
	
	avroJSONResults, err := pb.benchmarkAvroJSON("user")
	if err != nil {
		return fmt.Errorf("Avro JSON benchmark failed: %w", err)
	}

	avroBinaryResults, err := pb.benchmarkAvroBinary("user")
	if err != nil {
		return fmt.Errorf("Avro binary benchmark failed: %w", err)
	}

	stdJSONResults, err := pb.benchmarkStandardJSON("user")
	if err != nil {
		return fmt.Errorf("Standard JSON benchmark failed: %w", err)
	}

	// Display results
	pb.displayResults("User", []BenchmarkResults{avroJSONResults, avroBinaryResults, stdJSONResults})

	// Run product benchmarks
	fmt.Println("--- Product Serialization Benchmarks ---")
	
	avroJSONProductResults, err := pb.benchmarkAvroJSON("product")
	if err != nil {
		return fmt.Errorf("Avro JSON product benchmark failed: %w", err)
	}

	avroBinaryProductResults, err := pb.benchmarkAvroBinary("product")
	if err != nil {
		return fmt.Errorf("Avro binary product benchmark failed: %w", err)
	}

	stdJSONProductResults, err := pb.benchmarkStandardJSON("product")
	if err != nil {
		return fmt.Errorf("Standard JSON product benchmark failed: %w", err)
	}

	// Display results
	pb.displayResults("Product", []BenchmarkResults{avroJSONProductResults, avroBinaryProductResults, stdJSONProductResults})

	// Show summary
	pb.showSummary([]BenchmarkResults{avroJSONResults, avroBinaryResults, stdJSONResults})

	return nil
}

// benchmarkAvroJSON benchmarks Avro JSON serialization
func (pb *PerformanceBenchmark) benchmarkAvroJSON(dataType string) (BenchmarkResults, error) {
	var memBefore, memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()

	var totalSize int
	var iterations int

	if dataType == "user" {
		for _, user := range pb.users {
			data, err := pb.manager.SerializeUserJSON(user)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization 
			_, err = pb.manager.DeserializeUserJSON(data)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	} else {
		for _, product := range pb.products {
			data, err := pb.manager.SerializeProductJSON(product)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization
			_, err = pb.manager.DeserializeProductJSON(data)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	}

	elapsed := time.Since(startTime)
	runtime.ReadMemStats(&memAfter)

	return BenchmarkResults{
		Format:              "Avro JSON",
		SerializationTime:   elapsed / 2, // Approximate since we do both ser/deser
		DeserializationTime: elapsed / 2,
		SerializedSize:      totalSize / iterations,
		MemoryUsage:         int64(memAfter.TotalAlloc - memBefore.TotalAlloc),
		ItemsPerSecond:      float64(iterations*2) / elapsed.Seconds(), // ser + deser
	}, nil
}

// benchmarkAvroBinary benchmarks Avro binary serialization
func (pb *PerformanceBenchmark) benchmarkAvroBinary(dataType string) (BenchmarkResults, error) {
	var memBefore, memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()

	var totalSize int
	var iterations int

	if dataType == "user" {
		for _, user := range pb.users {
			data, err := pb.manager.SerializeUserBinary(user)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization
			_, err = pb.manager.DeserializeUserBinary(data)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	} else {
		for _, product := range pb.products {
			data, err := pb.manager.SerializeProductBinary(product)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization
			_, err = pb.manager.DeserializeProductBinary(data)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	}

	elapsed := time.Since(startTime)
	runtime.ReadMemStats(&memAfter)

	return BenchmarkResults{
		Format:              "Avro Binary",
		SerializationTime:   elapsed / 2,
		DeserializationTime: elapsed / 2,
		SerializedSize:      totalSize / iterations,
		MemoryUsage:         int64(memAfter.TotalAlloc - memBefore.TotalAlloc),
		ItemsPerSecond:      float64(iterations*2) / elapsed.Seconds(),
	}, nil
}

// benchmarkStandardJSON benchmarks standard Go JSON serialization
func (pb *PerformanceBenchmark) benchmarkStandardJSON(dataType string) (BenchmarkResults, error) {
	var memBefore, memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()

	var totalSize int
	var iterations int

	if dataType == "user" {
		for _, user := range pb.users {
			data, err := json.Marshal(user)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization
			var deserializedUser User
			err = json.Unmarshal(data, &deserializedUser)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	} else {
		for _, product := range pb.products {
			data, err := json.Marshal(product)
			if err != nil {
				return BenchmarkResults{}, err
			}
			totalSize += len(data)
			iterations++

			// Test deserialization
			var deserializedProduct Product
			err = json.Unmarshal(data, &deserializedProduct)
			if err != nil {
				return BenchmarkResults{}, err
			}
		}
	}

	elapsed := time.Since(startTime)
	runtime.ReadMemStats(&memAfter)

	return BenchmarkResults{
		Format:              "Standard JSON",
		SerializationTime:   elapsed / 2,
		DeserializationTime: elapsed / 2,
		SerializedSize:      totalSize / iterations,
		MemoryUsage:         int64(memAfter.TotalAlloc - memBefore.TotalAlloc),
		ItemsPerSecond:      float64(iterations*2) / elapsed.Seconds(),
	}, nil
}

// displayResults displays benchmark results in a formatted table
func (pb *PerformanceBenchmark) displayResults(dataType string, results []BenchmarkResults) {
	fmt.Printf("\n%s Serialization Performance:\n", dataType)
	fmt.Printf("%-15s %-12s %-15s %-12s %-15s %-12s\n", 
		"Format", "Ser Time", "Deser Time", "Size (B)", "Memory (KB)", "Items/sec")
	fmt.Printf("%-15s %-12s %-15s %-12s %-15s %-12s\n", 
		"------", "--------", "----------", "--------", "----------", "---------")

	for _, result := range results {
		fmt.Printf("%-15s %-12s %-15s %-12d %-15.1f %-12.0f\n",
			result.Format,
			formatDuration(result.SerializationTime),
			formatDuration(result.DeserializationTime),
			result.SerializedSize,
			float64(result.MemoryUsage)/1024,
			result.ItemsPerSecond)
	}

	// Show size comparisons
	if len(results) > 1 {
		fmt.Println("\nSize Comparison:")
		baseSize := results[0].SerializedSize
		for _, result := range results {
			if result.SerializedSize != baseSize {
				savings := float64(baseSize-result.SerializedSize) / float64(baseSize) * 100
				fmt.Printf("  %s vs %s: %.1f%% size difference\n", 
					results[0].Format, result.Format, savings)
			}
		}
	}
}

// showSummary displays an overall performance summary
func (pb *PerformanceBenchmark) showSummary(results []BenchmarkResults) {
	fmt.Println("\n=== Performance Summary ===")
	
	// Find fastest serializer
	fastest := results[0]
	for _, result := range results[1:] {
		if result.ItemsPerSecond > fastest.ItemsPerSecond {
			fastest = result
		}
	}
	fmt.Printf("✓ Fastest overall: %s (%.0f items/sec)\n", fastest.Format, fastest.ItemsPerSecond)

	// Find most memory efficient
	mostEfficient := results[0]
	for _, result := range results[1:] {
		if result.MemoryUsage < mostEfficient.MemoryUsage {
			mostEfficient = result
		}
	}
	fmt.Printf("✓ Most memory efficient: %s (%d KB)\n", 
		mostEfficient.Format, mostEfficient.MemoryUsage/1024)

	// Find smallest serialized size
	smallest := results[0]
	for _, result := range results[1:] {
		if result.SerializedSize < smallest.SerializedSize {
			smallest = result
		}
	}
	fmt.Printf("✓ Smallest serialized size: %s (%d bytes)\n", 
		smallest.Format, smallest.SerializedSize)

	fmt.Println("\nKey Findings:")
	fmt.Println("• Avro provides schema validation and evolution capabilities")
	fmt.Println("• Binary formats typically offer better compression")
	fmt.Println("• JSON formats are more human-readable and debuggable")
	fmt.Println("• Performance varies based on data structure complexity")
}

// formatDuration formats duration for display
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.1fμs", float64(d.Nanoseconds())/1000)
	} else if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
}

// RunPerformanceComparison runs the complete performance comparison
func RunPerformanceComparison() error {
	benchmark, err := NewPerformanceBenchmark()
	if err != nil {
		return fmt.Errorf("failed to create benchmark: %w", err)
	}

	return benchmark.RunBenchmarks()
}