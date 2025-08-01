package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deleonn/gopr/internal/models"
)

type AnthropicProvider struct {
	apiKey string
	model  string
}

func NewAnthropicProvider(config models.AnthropicConfig) *AnthropicProvider {
	model := config.Model

	if model == "" {
		model = "claude-3-sonnet-20240229"
	}

	return &AnthropicProvider{
		apiKey: config.APIKey,
		model:  model,
	}
}

func (a *AnthropicProvider) GetName() string {
	return "Anthropic"
}

func (a *AnthropicProvider) GenerateResponse(ctx context.Context, prompt string, temperature float64) (string, error) {
	requestBody := map[string]any{
		"model":       a.model,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"temperature": temperature,
		"max_tokens":  4000,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", a.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		return "", fmt.Errorf("invalid response format: no content")
	}

	firstContent, ok := content[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid response format: invalid content")
	}

	text, ok := firstContent["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: no text")
	}

	return text, nil
}

