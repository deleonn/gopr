package service

import (
	"fmt"

	"github.com/deleonn/gopr/internal/models"
)

// ProviderFactory creates LLM providers based on configuration
type ProviderFactory struct{}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

// CreateProvider creates a new LLM provider based on the configuration
func (f *ProviderFactory) CreateProvider(config models.Config) (models.LLMProvider, error) {
	switch config.Provider {
	case models.ProviderOllama:
		ollamaConfig := models.OllamaConfig{
			BaseURL: config.BaseURL,
			Model:   config.Model,
		}
		return NewOllamaProvider(ollamaConfig), nil

	case models.ProviderOpenAI:
		if config.APIKey == "" {
			return nil, fmt.Errorf("API key is required for OpenAI provider")
		}
		openAIConfig := models.OpenAIConfig{
			APIKey:  config.APIKey,
			Model:   config.Model,
			BaseURL: config.BaseURL,
		}
		return NewOpenAIProvider(openAIConfig), nil

	case models.ProviderAnthropic:
		if config.APIKey == "" {
			return nil, fmt.Errorf("API key is required for Anthropic provider")
		}
		anthropicConfig := models.AnthropicConfig{
			APIKey: config.APIKey,
			Model:  config.Model,
		}
		return NewAnthropicProvider(anthropicConfig), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
} 