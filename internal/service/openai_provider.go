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

type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
}

func NewOpenAIProvider(config models.OpenAIConfig) *OpenAIProvider {
	model := config.Model

	if model == "" {
		model = "gpt-4"
	}

	return &OpenAIProvider{
		apiKey: config.APIKey,
		model:  config.Model,
	}
}

func (o *OpenAIProvider) GetName() string {
	return "OpenAI"
}

func (o *OpenAIProvider) GenerateResponse(ctx context.Context, prompt string, temperature float64) (string, error) {
	requestBody := map[string]any{
		"model":       o.model,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"temperature": temperature,
		"max_tokens":  4000,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openapi.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)

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

	choices, ok := result["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format: no choices")
	}

	choice, ok := choices[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid response format: invalid choice")
	}

	message, ok := choice["message"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid response format: no message")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: no content")
	}

	return content, nil
}

