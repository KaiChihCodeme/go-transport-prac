package parquet

import (
	"os"
	"testing"
)

func TestETLWorkflow(t *testing.T) {
	testDir := "tmp/test_etl_workflow"
	pipeline := NewDataPipeline(testDir)
	defer pipeline.CleanupWorkflow()

	err := pipeline.RunETLWorkflow()
	if err != nil {
		t.Fatalf("ETL workflow failed: %v", err)
	}

	t.Log("✓ ETL workflow completed successfully")
}

func TestBatchProcessing(t *testing.T) {
	testDir := "tmp/test_batch_processing"
	pipeline := NewDataPipeline(testDir)
	defer pipeline.CleanupWorkflow()

	err := pipeline.RunBatchProcessing()
	if err != nil {
		t.Fatalf("Batch processing failed: %v", err)
	}

	t.Log("✓ Batch processing completed successfully")
}

func TestAnalyticsWorkflow(t *testing.T) {
	testDir := "tmp/test_analytics_workflow"
	pipeline := NewDataPipeline(testDir)  
	defer pipeline.CleanupWorkflow()

	err := pipeline.RunAnalyticsWorkflow()
	if err != nil {
		t.Fatalf("Analytics workflow failed: %v", err)
	}

	t.Log("✓ Analytics workflow completed successfully")
}

func TestDataQualityCalculation(t *testing.T) {
	pipeline := NewDataPipeline("tmp/test_quality")
	defer pipeline.CleanupWorkflow()

	// Test high quality user
	highQualityUser := User{
		ID:     1,
		Email:  "test@example.com",
		Name:   "Test User",
		Status: "active",
		Profile: &Profile{
			FirstName: "Test",
			LastName:  "User",
			Phone:     "+1-555-0123",
			Address: &Address{
				Country: "USA",
			},
		},
	}

	score := pipeline.calculateDataQuality(highQualityUser)
	if score < 0.8 {
		t.Errorf("Expected high quality score (>0.8), got %.2f", score)
	}
	t.Logf("High quality user score: %.2f", score)

	// Test low quality user
	lowQualityUser := User{
		ID:   2,
		Name: "Incomplete User",
		// Missing email, status, profile details
	}

	score = pipeline.calculateDataQuality(lowQualityUser)
	if score > 0.5 {
		t.Errorf("Expected low quality score (<0.5), got %.2f", score)
	}
	t.Logf("Low quality user score: %.2f", score)
}

func TestDataTransformation(t *testing.T) {
	pipeline := NewDataPipeline("tmp/test_transform")
	defer pipeline.CleanupWorkflow()

	// Create test data with various status formats
	rawUsers := []User{
		{ID: 1, Email: "test1@example.com", Name: "Test One", Status: "ACTIVE"},
		{ID: 2, Email: "test2@example.com", Name: "Test Two", Status: "Inactive"},
		{ID: 3, Email: "test3@example.com", Name: "Test Three", Status: "suspended"},
	}

	transformed, err := pipeline.transformUserData(rawUsers)
	if err != nil {
		t.Fatalf("Transformation failed: %v", err)
	}

	// Verify normalization
	expected := []string{"active", "inactive", "suspended"}
	for i, user := range transformed {
		if user.Status != expected[i] {
			t.Errorf("User %d: expected status %s, got %s", i, expected[i], user.Status)
		}
	}

	// Verify metadata was added
	for i, user := range transformed {
		if user.Profile == nil || user.Profile.Metadata == nil {
			t.Errorf("User %d: missing metadata", i)
			continue
		}
		if user.Profile.Metadata["transformed"] == "" {
			t.Errorf("User %d: missing transformed timestamp", i)
		}
		if user.Profile.Metadata["quality_score"] == "" {
			t.Errorf("User %d: missing quality score", i)
		}
	}

	t.Log("✓ Data transformation completed successfully")
}

func TestPhoneNormalization(t *testing.T) {
	pipeline := NewDataPipeline("tmp/test_phone")
	defer os.RemoveAll("tmp/test_phone")

	testCases := []struct {
		input    string
		expected string
	}{
		{"+1-555-0123", "+1-555-0123"}, // Already normalized
		{"555-0123", "+1-555-0123"},   // Add country code
		{"+44-20-1234", "+44-20-1234"}, // International, keep as is
		{"", ""},                       // Empty, keep as is
	}

	for _, tc := range testCases {
		result := pipeline.normalizePhoneNumber(tc.input)
		if result != tc.expected {
			t.Errorf("Phone normalization: input %s, expected %s, got %s", 
				tc.input, tc.expected, result)
		}
	}

	t.Log("✓ Phone normalization tests passed")
}

func TestNameSplitting(t *testing.T) {
	pipeline := NewDataPipeline("tmp/test_names")
	defer os.RemoveAll("tmp/test_names")

	testCases := []struct {
		input    string
		expected []string
	}{
		{"John Doe", []string{"John", "Doe"}},
		{"Alice", []string{"Alice"}},
		{"Mary Jane Watson", []string{"Mary", "Jane"}}, // Only first two parts
		{"", []string{}},
	}

	for _, tc := range testCases {
		result := pipeline.splitFullName(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("Name split: input %s, expected %v, got %v", 
				tc.input, tc.expected, result)
			continue
		}
		for i, part := range result {
			if part != tc.expected[i] {
				t.Errorf("Name split part %d: input %s, expected %s, got %s", 
					i, tc.input, tc.expected[i], part)
			}
		}
	}

	t.Log("✓ Name splitting tests passed")
}