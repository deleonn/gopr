package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deleonn/gopr/internal/models"
	"github.com/deleonn/gopr/internal/service"
)

func loadConfig() models.Config {
	config := models.Config{
		Temperature: 0.1,
	}

	// Try to load from .goprrc in current directory
	if configFile := ".goprrc"; fileExists(configFile) {
		loadConfigFromFile(configFile, &config)
	}

	// Try to load from ~/.goprrc
	if homeDir, err := os.UserHomeDir(); err == nil {
		if configFile := filepath.Join(homeDir, ".goprrc"); fileExists(configFile) {
			loadConfigFromFile(configFile, &config)
		}
	}

	return config
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadConfigFromFile(filename string, config *models.Config) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "provider":
			config.Provider = models.ProviderType(value)
		case "model":
			config.Model = value
		case "api_key":
			config.APIKey = value
		case "base_url":
			config.BaseURL = value
		case "temperature":
			if temp, err := strconv.ParseFloat(value, 64); err == nil {
				config.Temperature = temp
			}
		}
	}
}

func main() {
	config := loadConfig()

	// Parse command line flags (these override config file)
	var (
		provider    = flag.String("provider", string(config.Provider), "LLM provider (ollama, openai, anthropic, deepseek)")
		model       = flag.String("model", config.Model, "Model to use")
		apiKey      = flag.String("api-key", config.APIKey, "API key for the provider")
		baseURL     = flag.String("base-url", config.BaseURL, "Base URL for the provider")
		temperature = flag.Float64("temperature", config.Temperature, "Temperature for generation")
		branch      = flag.String("branch", "main", "Branch for diff comparison")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
	)
	flag.Parse()

	// Override config with command line flags
	config.Provider = models.ProviderType(*provider)
	config.Model = *model
	config.APIKey = *apiKey
	config.BaseURL = *baseURL
	config.Temperature = *temperature

	prService, err := service.NewPRService(config, *branch)
	if err != nil {
		log.Fatalf("Failed to create PR service: %v", err)
	}

	description, err := prService.GeneratePRDescriptionFromBranch(*verbose)
	if err != nil {
		log.Fatalf("Failed to generate PR description: %v", err)
	}

	// Output the description to stdout (can be piped to gh or clipboard)
	fmt.Print(description)
}
