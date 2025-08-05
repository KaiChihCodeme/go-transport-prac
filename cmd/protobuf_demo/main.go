package main

import (
	"log"

	"go-transport-prac/pkg/sdl/protobuf"
)

func main() {
	examples := protobuf.NewExamples()
	
	if err := examples.RunAllExamples(); err != nil {
		log.Fatalf("Failed to run examples: %v", err)
	}
}