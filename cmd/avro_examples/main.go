package main

import (
	"log"
	"os"

	"go-transport-prac/pkg/sdl/avro"
)

func main() {
	// Create examples instance
	examples, err := avro.NewExamples()
	if err != nil {
		log.Fatalf("Failed to create examples: %v", err)
	}

	// Run all examples
	err = examples.RunAllExamples()
	if err != nil {
		log.Fatalf("Examples failed: %v", err)
	}

	// Cleanup
	err = examples.CleanupExamples()
	if err != nil {
		log.Printf("Cleanup warning: %v", err)
	}

	// Clean up temp directories
	os.RemoveAll("tmp")
}