package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/deleonn/gopr/internal/service"
)

type Config struct {
	OllamaURL string
	Model     string
}

func loadConfig() Config {
	config := Config{
		OllamaURL: "http://msi_th.home:11434",
		Model:     "qwen2.5-coder:14b-instruct-q8_0",
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

func loadConfigFromFile(filename string, config *Config) {
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
		case "ollama_url":
			config.OllamaURL = value
		case "model":
			config.Model = value
		}
	}
}

func main() {
	config := loadConfig()

	// Parse command line flags (these override config file)
	var (
		ollamaURL = flag.String("ollama-url", config.OllamaURL, "Ollama server URL")
		model     = flag.String("model", config.Model, "Ollama model to use")
		branch    = flag.String("branch", "main", "Branch for diff comparison")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
	)
	flag.Parse()

	prService := service.NewPRService(*ollamaURL, *model, *branch)

	description, err := prService.GeneratePRDescriptionFromBranch(*verbose)
	if err != nil {
		log.Fatalf("Failed to generate PR description: %v", err)
	}

	// Output the description to stdout (can be piped to gh or clipboard)
	fmt.Print(description)
}
