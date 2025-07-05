package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type PRService struct {
	ollamaURL string
	model     string
}

func NewPRService(ollamaURL, model string) *PRService {
	return &PRService{
		ollamaURL: ollamaURL,
		model:     model,
	}
}

// GeneratePRDescriptionFromBranch generates a PR description by comparing current branch with main
func (s *PRService) GeneratePRDescriptionFromBranch(verbose bool) (string, error) {
	// Get the current branch name
	currentBranch, err := s.getCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Current branch: %s\n", currentBranch)
	}

	// Get the diff between current branch and main
	diff, err := s.getBranchDiff()
	if err != nil {
		return "", fmt.Errorf("failed to get branch diff: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Diff length: %d characters\n", len(diff))
	}

	// Get commit messages since main
	commits, err := s.getCommitsSinceMain()
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Number of commits: %d\n", len(commits))
	}

	// Format the information for the LLM
	prompt := s.formatForLLM(currentBranch, diff, commits)

	// Generate description using Ollama
	description, err := s.callOllama(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate description: %w", err)
	}

	return description, nil
}

// getCurrentBranch gets the name of the current branch
func (s *PRService) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getBranchDiff gets the diff between current branch and main
func (s *PRService) getBranchDiff() (string, error) {
	cmd := exec.Command("git", "diff", "main...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %v", err)
	}
	return string(output), nil
}

// getCommitsSinceMain gets commit messages since the main branch
func (s *PRService) getCommitsSinceMain() ([]string, error) {
	cmd := exec.Command("git", "log", "main..", "--oneline", "--no-merges")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

// formatForLLM formats the information for optimal LLM input
func (s *PRService) formatForLLM(branchName, diff string, commits []string) string {
	var prompt strings.Builder
	
	prompt.WriteString("You are a helpful assistant that generates clear and concise Pull Request descriptions. ")
	prompt.WriteString("Based on the following information, generate a professional PR description in markdown format.\n\n")
	
	prompt.WriteString("## Context\n")
	prompt.WriteString(fmt.Sprintf("- **Branch**: %s\n", branchName))
	prompt.WriteString(fmt.Sprintf("- **Number of commits**: %d\n", len(commits)))
	
	if len(commits) > 0 {
		prompt.WriteString("\n## Commits\n")
		for _, commit := range commits {
			prompt.WriteString(fmt.Sprintf("- %s\n", commit))
		}
	}
	
	prompt.WriteString("\n## Changes\n")
	prompt.WriteString("```diff\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n```\n\n")
	
	prompt.WriteString("## Instructions\n")
	prompt.WriteString("Generate a PR description that includes:\n")
	prompt.WriteString("1. A brief summary of what this PR accomplishes\n")
	prompt.WriteString("2. Key changes made\n")
	prompt.WriteString("3. Any breaking changes or important notes\n")
	prompt.WriteString("4. Testing information if applicable\n\n")
	prompt.WriteString("Keep it professional, clear, and concise. Use markdown formatting.\n\n")
	prompt.WriteString("## PR Description\n")
	
	return prompt.String()
}

// callOllama makes a request to the Ollama API
func (s *PRService) callOllama(prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"model":  s.model,
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	ctx := context.Background()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.ollamaURL+"/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
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

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	return response, nil
}

// Legacy method for backward compatibility (can be removed if not needed)
func (s *PRService) GeneratePRDescription(req interface{}) (string, error) {
	return "", fmt.Errorf("this method is deprecated, use GeneratePRDescriptionFromBranch instead")
}
