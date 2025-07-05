package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/joeldeleon/pr_description_generator/internal/models"
	"golang.org/x/oauth2"
)

type PRService struct {
	apiKey string
	baseURL string
}

func NewPRService(apiKey, baseURL string) *PRService {
	return &PRService{
		apiKey: apiKey,
		baseURL: baseURL,
	}
}

func (s *PRService) GeneratePRDescription(req models.PRRequest) (string, error) {
	prompt := fmt.Sprintf(
		"Generate a PR description for the following changes:\n\nTitle: %s\nDescription: %s\nChanges:\n",
		req.Title, req.Description)

	for _, change := range req.Changes {
		prompt += fmt.Sprintf("- %s: %s\n", change.ChangeType, change.FilePath)
	}

	requestBody := map[string]interface{}{
		"model": "llm-model-name", // Replace with the actual model name
		"prompt": prompt,
		"max_tokens": 150,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/v1/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	text, ok := choices[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text format in response")
	}

	return text, nil
}

