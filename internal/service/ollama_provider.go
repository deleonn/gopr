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

type OllamaProvider struct {
	baseURL string
	model   string
}

func NewOllamaProvider(config models.OllamaConfig) *OllamaProvider {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := config.Model
	if model == "" {
		model = "qwen2.5-coder:14b-instruct-q8_0"
	}

	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
	}
}

func (o *OllamaProvider) GetName() string {
	return "Ollama"
}

func (o *OllamaProvider) GenerateResponse(ctx context.Context, prompt string, temperature float64) (string, error) {
	requestBody := map[string]any{
		"model":       o.model,
		"prompt":      prompt,
		"stream":      false,
		"temperature": temperature,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

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

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	return response, nil
}

