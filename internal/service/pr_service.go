package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/deleonn/gopr/internal/models"
)

type PRService struct {
	provider models.LLMProvider
	branch   string
}

func NewPRService(config models.Config, branch string) (*PRService, error) {
	factory := NewProviderFactory()
	provider, err := factory.CreateProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	return &PRService{
		provider: provider,
		branch:   branch,
	}, nil
}

// GeneratePRDescriptionFromBranch generates a PR description by comparing current branch with the provided branch
func (s *PRService) GeneratePRDescriptionFromBranch(verbose bool) (string, error) {
	// Get the current branch name
	currentBranch, err := s.getCurrentBranch()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Current provider: %s\n", s.provider.GetName())
		fmt.Fprintf(os.Stderr, "Current model: %s\n", s.provider.GetModel())
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Current branch: %s\n", currentBranch)
	}

	// Get the diff between current branch and the selected one
	diff, err := s.getBranchDiff(s.branch)
	if err != nil {
		return "", fmt.Errorf("failed to get branch diff: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Diff length: %d characters\n", len(diff))
	}

	// Get commit messages since branch
	commits, err := s.getCommitsSince()
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Number of commits: %d\n", len(commits))
	}

	// Analyze file types for better context
	fileAnalysis := s.analyzeFileTypes(diff)
	if verbose {
		fmt.Fprintf(os.Stderr, "File analysis: %s\n", fileAnalysis)
	}

	// Format the information for the LLM
	prompt := s.formatForLLM(currentBranch, diff, commits, fileAnalysis)

	// Generate description using LLM provider with retry logic
	var description string
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if verbose && attempt > 1 {
			fmt.Fprintf(os.Stderr, "Retry attempt %d/%d\n", attempt, maxRetries)
		}

		description, err = s.callLLMProvider(prompt)
		if err != nil {
			if attempt == maxRetries {
				return "", fmt.Errorf("failed to generate description after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		// Validate the response
		if s.validateResponse(description) {
			break
		} else if attempt == maxRetries {
			if verbose {
				fmt.Fprintf(os.Stderr, "Warning: Generated response may be too generic\n")
			}
		} else {
			if verbose {
				fmt.Fprintf(os.Stderr, "Response too generic, retrying...\n")
			}
			time.Sleep(time.Duration(attempt) * time.Second)
		}
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

// getBranchDiff gets the diff between current branch and the provided branch
func (s *PRService) getBranchDiff(branch string) (string, error) {
	cmd := exec.Command("git", "diff", fmt.Sprintf("%s...", branch))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %v", err)
	}
	return string(output), nil
}

// getCommitsSince gets commit messages since the provided branch
func (s *PRService) getCommitsSince() ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("%s..", s.branch), "--oneline", "--no-merges")
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

// analyzeFileTypes analyzes the diff to understand what types of files were changed
func (s *PRService) analyzeFileTypes(diff string) string {
	var analysis strings.Builder
	analysis.WriteString("## File Analysis\n")

	// Count different file types
	fileTypes := make(map[string]int)
	lines := strings.Split(diff, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") || strings.HasPrefix(line, "--- a/") {
			parts := strings.Split(line, "/")
			if len(parts) > 1 {
				fileName := parts[len(parts)-1]
				ext := ""
				if idx := strings.LastIndex(fileName, "."); idx != -1 {
					ext = fileName[idx:]
				}
				if ext != "" {
					fileTypes[ext]++
				}
			}
		}
	}

	if len(fileTypes) > 0 {
		analysis.WriteString("Files changed by type:\n")
		for ext, count := range fileTypes {
			analysis.WriteString(fmt.Sprintf("- %s: %d files\n", ext, count))
		}
	} else {
		analysis.WriteString("No file type analysis available\n")
	}

	return analysis.String()
}

// formatForLLM formats the information for optimal LLM input
func (s *PRService) formatForLLM(branchName, diff string, commits []string, fileAnalysis string) string {
	var prompt strings.Builder

	prompt.WriteString("You are analyzing a Git repository to generate an accurate PR description. ")
	prompt.WriteString("Base your response ONLY on the actual code changes shown below. ")
	prompt.WriteString("Do NOT make assumptions or generic statements. ")
	prompt.WriteString("If the changes are unclear, be specific about what you can see.\n\n")

	prompt.WriteString("## Repository Context\n")
	prompt.WriteString(fmt.Sprintf("Current branch: %s\n", branchName))
	prompt.WriteString(fmt.Sprintf("Number of commits since %s: %d\n", s.branch, len(commits)))

	if len(commits) > 0 {
		prompt.WriteString("\nCommit messages:\n")
		for _, commit := range commits {
			prompt.WriteString(fmt.Sprintf("- %s\n", commit))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString(fileAnalysis)
	prompt.WriteString("\n")

	prompt.WriteString("## Actual Code Changes (git diff)\n")
	if diff == "" {
		prompt.WriteString("No code changes detected (empty diff)\n\n")
	} else {
		prompt.WriteString("```diff\n")
		prompt.WriteString(diff)
		prompt.WriteString("\n```\n\n")
	}

	prompt.WriteString("## Instructions\n")
	prompt.WriteString("Analyze the code changes above and generate a PR description. ")
	prompt.WriteString("Be specific about what files were changed and what functionality was added/modified/removed. ")
	prompt.WriteString("If you cannot determine the purpose from the code, say so clearly.\n\n")
	prompt.WriteString("Respond with ONLY the PR description in this exact format:\n\n")
	prompt.WriteString("# TL;DR\n")
	prompt.WriteString("[Specific summary based on actual changes]\n\n")
	prompt.WriteString("# What's changed?\n")
	prompt.WriteString("- [Specific change based on diff]\n")
	prompt.WriteString("- [Another specific change]\n\n")
	prompt.WriteString("# How to test?\n")
	prompt.WriteString("1. [Specific test step related to changes]\n")
	prompt.WriteString("2. [Another specific test step]\n\n")
	prompt.WriteString("# Why make this change?\n")
	prompt.WriteString("[Reasoning based on actual code changes]\n\n")
	prompt.WriteString("# Breaking changes or important notes\n")
	prompt.WriteString("- [Important note based on actual changes]\n")
	prompt.WriteString("- [Another important note if applicable]\n\n")
	prompt.WriteString("## Your Response\n")

	return prompt.String()
}

// validateResponse checks if the response is too generic
func (s *PRService) validateResponse(response string) bool {
	genericPhrases := []string{
		"improvements to the codebase",
		"enhancing user experience",
		"fixing minor bugs",
		"improved functionality",
		"better performance",
		"general improvements",
		"enhancing the overall",
		"improving the system",
		"better user experience",
		"enhanced features",
		"improved performance",
		"better functionality",
	}

	responseLower := strings.ToLower(response)
	for _, phrase := range genericPhrases {
		if strings.Contains(responseLower, phrase) {
			return false // Too generic
		}
	}
	return true
}

// callLLMProvider makes a request to the configured LLM provider
func (s *PRService) callLLMProvider(prompt string) (string, error) {
	ctx := context.Background()
	return s.provider.GenerateResponse(ctx, prompt, 0.1) // Low temperature for more focused responses
}
