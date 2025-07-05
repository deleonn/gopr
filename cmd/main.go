package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/deleonn/gopr/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Parse command line flags
	var (
		ollamaURL = flag.String("ollama-url", "http://localhost:11434", "Ollama server URL")
		model     = flag.String("model", "llama3.2", "Ollama model to use")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
	)
	flag.Parse()

	// Initialize the PR service
	prService := service.NewPRService(*ollamaURL, *model)

	// Generate PR description from branch comparison
	description, err := prService.GeneratePRDescriptionFromBranch(*verbose)
	if err != nil {
		log.Fatalf("Failed to generate PR description: %v", err)
	}

	// Output the description to stdout (can be piped to gh or clipboard)
	fmt.Print(description)
}
