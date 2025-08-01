package models

import (
	"context"
)

// LLMProvider defines the interface for different LLM providers
type LLMProvider interface {
	GenerateResponse(ctx context.Context, prompt string, temperature float64) (string, error)
	GetName() string
}

// ProviderType represents the type of LLM provider
type ProviderType string

const (
	ProviderOllama   ProviderType = "ollama"
	ProviderOpenAI   ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
)

// Config holds the configuration for the application
type Config struct {
	Provider   ProviderType `json:"provider"`
	Model      string       `json:"model"`
	APIKey     string       `json:"api_key,omitempty"`
	BaseURL    string       `json:"base_url,omitempty"`
	Temperature float64      `json:"temperature"`
}

// OllamaConfig holds Ollama-specific configuration
type OllamaConfig struct {
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

// OpenAIConfig holds OpenAI-specific configuration
type OpenAIConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
	BaseURL string `json:"base_url,omitempty"`
}

// AnthropicConfig holds Anthropic-specific configuration
type AnthropicConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
} 