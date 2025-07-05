package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/joeldeleon/pr_description_generator/internal/api"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Initialize the API
	apiServer := api.NewAPI()

	// Start the server
	log.Println("Starting server on port 8080")
	if err := apiServer.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

