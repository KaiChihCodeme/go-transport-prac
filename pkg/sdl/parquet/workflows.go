package parquet

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// DataPipeline demonstrates a complete data processing workflow using Parquet
type DataPipeline struct {
	manager     *SimpleManager
	inputDir    string
	outputDir   string
	processedDir string
}

// NewDataPipeline creates a new data processing pipeline
func NewDataPipeline(baseDir string) *DataPipeline {
	return &DataPipeline{
		manager:      NewSimpleManager(filepath.Join(baseDir, "data")),
		inputDir:     filepath.Join(baseDir, "input"),
		outputDir:    filepath.Join(baseDir, "output"),
		processedDir: filepath.Join(baseDir, "processed"),
	}
}

// RunETLWorkflow demonstrates an ETL (Extract, Transform, Load) workflow
func (dp *DataPipeline) RunETLWorkflow() error {
	fmt.Println("=== ETL Workflow with Parquet ===")
	
	// 1. Extract: Generate sample data (simulating data extraction)
	rawUsers, err := dp.extractUserData()
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}
	fmt.Printf("✓ Extracted %d user records\n", len(rawUsers))
	
	// 2. Transform: Clean and enhance the data
	transformedUsers, err := dp.transformUserData(rawUsers)
	if err != nil {
		return fmt.Errorf("transformation failed: %w", err)
	}
	fmt.Printf("✓ Transformed %d user records\n", len(transformedUsers))
	
	// 3. Load: Save to Parquet format
	if err := dp.loadUserData(transformedUsers); err != nil {
		return fmt.Errorf("loading failed: %w", err)
	}
	fmt.Printf("✓ Loaded data to Parquet format\n")
	
	// 4. Verify: Read back and validate
	if err := dp.verifyLoadedData(); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Printf("✓ Data verification successful\n")
	
	return nil
}

// extractUserData simulates extracting data from various sources
func (dp *DataPipeline) extractUserData() ([]User, error) {
	// Simulate data from different sources with varying quality
	rawData := []struct {
		id       int64
		email    string
		name     string
		status   string
		phone    string
		city     string
		country  string
	}{
		{1, "alice@example.com", "Alice Smith", "active", "+1-555-0001", "New York", "USA"},
		{2, "bob@test.com", "Bob Johnson", "ACTIVE", "555-0002", "San Francisco", "USA"},
		{3, "charlie@demo.org", "", "inactive", "+44-20-1234", "London", "UK"},
		{4, "diana@sample.net", "Diana Prince", "suspended", "", "Toronto", "Canada"},
		{5, "eve@example.co.uk", "Eve Wilson", "Active", "+33-1-4567", "Paris", "France"},
	}
	
	users := make([]User, len(rawData))
	now := time.Now()
	
	for i, raw := range rawData {
		// Convert raw data to User struct (minimal transformation here)
		name := raw.name
		if name == "" {
			name = fmt.Sprintf("User %d", raw.id)
		}
		
		users[i] = User{
			ID:     raw.id,
			Email:  raw.email,
			Name:   name,
			Status: raw.status, // Will be normalized in transform step
			Profile: &Profile{
				Phone: raw.phone,
				Address: &Address{
					City:    raw.city,
					Country: raw.country,
				},
				Metadata: map[string]string{
					"source":    "raw_extraction",
					"extracted": now.Format(time.RFC3339),
				},
			},
			CreatedAt: now.Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt: now,
		}
	}
	
	return users, nil
}

// transformUserData cleans and enhances the extracted data
func (dp *DataPipeline) transformUserData(users []User) ([]User, error) {
	fmt.Println("Applying data transformations...")
	
	transformed := make([]User, len(users))
	
	for i, user := range users {
		// Copy the user
		transformed[i] = user
		
		// 1. Normalize status values
		switch user.Status {
		case "ACTIVE", "Active", "active":
			transformed[i].Status = "active"
		case "INACTIVE", "Inactive", "inactive":
			transformed[i].Status = "inactive"
		case "SUSPENDED", "Suspended", "suspended":
			transformed[i].Status = "suspended"
		default:
			transformed[i].Status = "unknown"
		}
		
		// 2. Normalize phone numbers
		if user.Profile != nil && user.Profile.Phone != "" {
			transformed[i].Profile.Phone = dp.normalizePhoneNumber(user.Profile.Phone)
		}
		
		// 3. Add computed fields
		if transformed[i].Profile == nil {
			transformed[i].Profile = &Profile{}
		}
		
		if transformed[i].Profile.Metadata == nil {
			transformed[i].Profile.Metadata = make(map[string]string)
		}
		
		// Add transformation metadata
		transformed[i].Profile.Metadata["transformed"] = time.Now().Format(time.RFC3339)
		transformed[i].Profile.Metadata["status_normalized"] = "true"
		
		// 4. Extract name parts if available
		if transformed[i].Profile.FirstName == "" && transformed[i].Name != "" {
			parts := dp.splitFullName(transformed[i].Name)
			transformed[i].Profile.FirstName = parts[0]
			if len(parts) > 1 {
				transformed[i].Profile.LastName = parts[1]
			}
		}
		
		// 5. Add data quality scores
		qualityScore := dp.calculateDataQuality(transformed[i])
		transformed[i].Profile.Metadata["quality_score"] = fmt.Sprintf("%.2f", qualityScore)
	}
	
	fmt.Printf("  - Normalized %d status values\n", len(transformed))
	fmt.Printf("  - Enhanced %d user profiles\n", len(transformed))
	
	return transformed, nil
}

// normalizePhoneNumber normalizes phone number format
func (dp *DataPipeline) normalizePhoneNumber(phone string) string {
	// Simple normalization - in real world this would be more sophisticated
	if len(phone) > 0 && phone[0] != '+' {
		// Add country code for US numbers
		if len(phone) == 8 || (len(phone) == 12 && phone[:3] == "555") {
			return "+1-" + phone
		}
	}
	return phone
}

// splitFullName splits full name into parts
func (dp *DataPipeline) splitFullName(fullName string) []string {
	// Simple split - real implementation would handle edge cases
	parts := []string{}
	if fullName != "" {
		// Split on space and take first two parts
		words := []string{}
		word := ""
		for _, r := range fullName {
			if r == ' ' {
				if word != "" {
					words = append(words, word)
					word = ""
				}
			} else {
				word += string(r)
			}
		}
		if word != "" {
			words = append(words, word)
		}
		
		if len(words) > 0 {
			parts = append(parts, words[0])
			if len(words) > 1 {
				parts = append(parts, words[1])
			}
		}
	}
	return parts
}

// calculateDataQuality calculates a data quality score (0-1)
func (dp *DataPipeline) calculateDataQuality(user User) float64 {
	score := 0.0
	maxScore := 10.0
	
	// Check required fields
	if user.ID > 0 {
		score += 2.0
	}
	if user.Email != "" {
		score += 2.0
	}
	if user.Name != "" {
		score += 1.0
	}
	if user.Status != "unknown" {
		score += 1.0
	}
	
	// Check profile completeness
	if user.Profile != nil {
		if user.Profile.FirstName != "" {
			score += 1.0
		}
		if user.Profile.LastName != "" {
			score += 1.0
		}
		if user.Profile.Phone != "" {
			score += 1.0
		}
		if user.Profile.Address != nil && user.Profile.Address.Country != "" {
			score += 1.0
		}
	}
	
	return score / maxScore
}

// loadUserData saves transformed data to Parquet
func (dp *DataPipeline) loadUserData(users []User) error {
	// Create output directory
	if err := os.MkdirAll(dp.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Save to Parquet with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("users_processed_%s.parquet", timestamp)
	
	outputManager := NewSimpleManager(dp.outputDir)
	return outputManager.WriteUsers(filename, users)
}

// verifyLoadedData reads back and validates the loaded data
func (dp *DataPipeline) verifyLoadedData() error {
	outputManager := NewSimpleManager(dp.outputDir)
	outputFiles, err := outputManager.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list output files: %w", err)
	}
	
	if len(outputFiles) == 0 {
		return fmt.Errorf("no output files found")
	}
	
	// Verify the most recent file
	latestFile := outputFiles[len(outputFiles)-1]
	users, err := outputManager.ReadUsers(latestFile)
	if err != nil {
		return fmt.Errorf("failed to read back data: %w", err)
	}
	
	// Validate data quality
	totalQuality := 0.0
	for _, user := range users {
		quality := dp.calculateDataQuality(user)
		totalQuality += quality
	}
	
	avgQuality := totalQuality / float64(len(users))
	fmt.Printf("  - Validated %d records\n", len(users))
	fmt.Printf("  - Average data quality: %.2f\n", avgQuality)
	
	if avgQuality < 0.7 {
		return fmt.Errorf("data quality too low: %.2f < 0.7", avgQuality)
	}
	
	return nil
}

// RunBatchProcessing demonstrates batch processing workflow
func (dp *DataPipeline) RunBatchProcessing() error {
	fmt.Println("=== Batch Processing Workflow ===")
	
	// Create multiple batches of data
	batchSize := 1000
	numBatches := 5
	
	fmt.Printf("Processing %d batches of %d records each...\n", numBatches, batchSize)
	
	for batch := 0; batch < numBatches; batch++ {
		// Generate batch data
		users := dp.generateBatchData(batch, batchSize)
		
		// Process batch
		filename := fmt.Sprintf("batch_%03d.parquet", batch)
		if err := dp.manager.WriteUsers(filename, users); err != nil {
			return fmt.Errorf("failed to write batch %d: %w", batch, err)
		}
		
		fmt.Printf("  ✓ Processed batch %d: %d records\n", batch, len(users))
	}
	
	// Aggregate results
	return dp.aggregateBatches()
}

// generateBatchData creates sample data for batch processing
func (dp *DataPipeline) generateBatchData(batchNum, size int) []User {
	users := make([]User, size)
	baseTime := time.Now().Add(-time.Duration(batchNum*24) * time.Hour)
	
	for i := 0; i < size; i++ {
		userID := int64(batchNum*size + i + 1)
		users[i] = User{
			ID:     userID,
			Email:  fmt.Sprintf("batch%d_user%d@example.com", batchNum, i),
			Name:   fmt.Sprintf("Batch %d User %d", batchNum, i),
			Status: []string{"active", "inactive", "suspended"}[i%3],
			Profile: &Profile{
				FirstName: fmt.Sprintf("First%d", i),
				LastName:  fmt.Sprintf("Last%d", i),
				Phone:     fmt.Sprintf("+1-555-%04d", i%10000),
				Address: &Address{
					City:    fmt.Sprintf("City%d", i%100),
					Country: []string{"USA", "Canada", "UK", "France", "Germany"}[i%5],
				},
				Interests: []string{
					fmt.Sprintf("interest%d", i%10),
					fmt.Sprintf("hobby%d", i%5),
				},
				Metadata: map[string]string{
					"batch":     fmt.Sprintf("%d", batchNum),
					"batch_pos": fmt.Sprintf("%d", i),
					"generated": baseTime.Format(time.RFC3339),
				},
			},
			CreatedAt: baseTime.Add(time.Duration(i) * time.Minute),
			UpdatedAt: time.Now(),
		}
	}
	
	return users
}

// aggregateBatches combines all batch files into summary statistics
func (dp *DataPipeline) aggregateBatches() error {
	fmt.Println("Aggregating batch results...")
	
	files, err := dp.manager.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}
	
	totalUsers := 0
	statusCounts := make(map[string]int)
	countryCounts := make(map[string]int)
	
	for _, filename := range files {
		if len(filename) > 5 && filename[:5] == "batch" {
			users, err := dp.manager.ReadUsers(filename)
			if err != nil {
				log.Printf("Warning: failed to read %s: %v", filename, err)
				continue
			}
			
			totalUsers += len(users)
			
			// Aggregate statistics
			for _, user := range users {
				statusCounts[user.Status]++
				if user.Profile != nil && user.Profile.Address != nil {
					countryCounts[user.Profile.Address.Country]++
				}
			}
		}
	}
	
	fmt.Printf("✓ Aggregation complete:\n")
	fmt.Printf("  - Total users processed: %d\n", totalUsers)
	fmt.Printf("  - Status distribution:\n")
	for status, count := range statusCounts {
		fmt.Printf("    %s: %d\n", status, count)
	}
	fmt.Printf("  - Country distribution:\n")
	for country, count := range countryCounts {
		fmt.Printf("    %s: %d\n", country, count)
	}
	
	return nil
}

// CleanupWorkflow removes all generated files
func (dp *DataPipeline) CleanupWorkflow() error {
	fmt.Println("=== Cleaning up workflow files ===")
	
	dirs := []string{
		dp.manager.baseDir,
		dp.inputDir,
		dp.outputDir,
		dp.processedDir,
	}
	
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("Warning: failed to remove %s: %v", dir, err)
		} else {
			fmt.Printf("✓ Removed %s\n", dir)
		}
	}
	
	return nil
}

// RunAnalyticsWorkflow demonstrates analytics data processing
func (dp *DataPipeline) RunAnalyticsWorkflow() error {
	fmt.Println("=== Analytics Workflow ===")
	
	// Generate time-series analytics data
	analyticsData := dp.generateAnalyticsData(24, 100) // 24 hours, 100 events per hour
	
	// Save analytics data
	filename := "analytics_data.parquet"
	if err := dp.writeAnalyticsData(filename, analyticsData); err != nil {
		return fmt.Errorf("failed to save analytics data: %w", err)
	}
	
	fmt.Printf("✓ Generated %d analytics events\n", len(analyticsData))
	
	// Process analytics data
	return dp.processAnalyticsData(filename)
}

// generateAnalyticsData creates sample analytics events
func (dp *DataPipeline) generateAnalyticsData(hours, eventsPerHour int) []Analytics {
	totalEvents := hours * eventsPerHour
	events := make([]Analytics, totalEvents)
	
	baseTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	eventTypes := []string{"page_view", "click", "purchase", "signup", "logout"}
	platforms := []string{"web", "mobile", "desktop"}
	countries := []string{"US", "CA", "GB", "DE", "FR", "JP", "AU"}
	
	for i := 0; i < totalEvents; i++ {
		hour := i / eventsPerHour
		eventTime := baseTime.Add(time.Duration(hour)*time.Hour + time.Duration(i%eventsPerHour)*time.Minute)
		
		events[i] = Analytics{
			ID:        int64(i + 1),
			EventType: eventTypes[i%len(eventTypes)],
			UserID:    int64((i % 1000) + 1),
			SessionID: fmt.Sprintf("session_%d", i%50),
			Timestamp: eventTime,
			Properties: map[string]string{
				"page":     fmt.Sprintf("/page/%d", i%10),
				"source":   "organic",
				"campaign": fmt.Sprintf("camp_%d", i%5),
			},
			Metrics: map[string]float64{
				"duration": float64(i%300 + 30),
				"value":    float64(i%100),
				"score":    float64(i%10) / 10.0,
			},
			DeviceInfo: &DeviceInfo{
				Platform: platforms[i%len(platforms)],
				Browser:  "chrome",
				Mobile:   platforms[i%len(platforms)] == "mobile",
			},
			Location: &Location{
				Country: countries[i%len(countries)],
				City:    fmt.Sprintf("City%d", i%20),
			},
		}
	}
	
	return events
}

// writeAnalyticsData saves analytics data (simplified version without full manager)
func (dp *DataPipeline) writeAnalyticsData(filename string, data []Analytics) error {
	// This is a simplified implementation - in full version we'd use the complete manager
	fmt.Printf("Writing %d analytics events to %s\n", len(data), filename)
	return nil
}

// processAnalyticsData analyzes the analytics data
func (dp *DataPipeline) processAnalyticsData(filename string) error {
	fmt.Println("Processing analytics data...")
	
	// Simulate analytics processing
	fmt.Println("  ✓ Calculated conversion rates")
	fmt.Println("  ✓ Generated user segments")
	fmt.Println("  ✓ Computed engagement metrics")
	fmt.Println("  ✓ Created daily/hourly aggregations")
	
	return nil
}