package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

type ReviewRules struct {
	Rules []string `yaml:"rules"`
}

func main() {
	ctx := context.Background()

	// Get environment variables
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	githubToken := os.Getenv("GITHUB_TOKEN")
	repoFullName := os.Getenv("GITHUB_REPOSITORY")
	prNumber := os.Getenv("PR_NUMBER")

	if apiKey == "" || githubToken == "" || repoFullName == "" || prNumber == "" {
		log.Fatal("Missing required environment variables: ANTHROPIC_API_KEY, GITHUB_TOKEN, GITHUB_REPOSITORY, PR_NUMBER")
	}

	// Parse repository owner and name
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		log.Fatalf("Invalid repository format: %s", repoFullName)
	}
	owner, repo := parts[0], parts[1]

	// Setup GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// Get PR number
	var prNum int
	fmt.Sscanf(prNumber, "%d", &prNum)

	// Get PR diff
	diff, err := getPRDiff(ctx, ghClient, owner, repo, prNum)
	if err != nil {
		log.Fatalf("Failed to get PR diff: %v", err)
	}

	if diff == "" {
		fmt.Println("No changes found in PR")
		return
	}

	// Read copilot instructions
	rules, err := readCopilotInstructions(ctx, ghClient, owner, repo)
	if err != nil {
		log.Printf("Warning: Could not read .github/copilot-instructions.yml: %v", err)
		rules = []string{"Check for code quality issues", "Look for potential bugs", "Suggest improvements"}
	}

	// Perform code review using Anthropic
	review, err := performCodeReview(ctx, apiKey, diff, rules)
	if err != nil {
		log.Fatalf("Failed to perform code review: %v", err)
	}

	// Post review as PR comment
	if err := postReviewComment(ctx, ghClient, owner, repo, prNum, review); err != nil {
		log.Fatalf("Failed to post review comment: %v", err)
	}

	fmt.Println("Code review completed successfully")
}

func getPRDiff(ctx context.Context, client *github.Client, owner, repo string, prNum int) (string, error) {
	// Get PR files
	opts := &github.ListOptions{PerPage: 100}
	var allFiles []*github.CommitFile

	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, prNum, opts)
		if err != nil {
			return "", err
		}
		allFiles = append(allFiles, files...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Build diff string
	var diffBuilder strings.Builder
	for _, file := range allFiles {
		if file.Patch != nil {
			diffBuilder.WriteString(fmt.Sprintf("\n=== %s ===\n", file.GetFilename()))
			diffBuilder.WriteString(*file.Patch)
			diffBuilder.WriteString("\n")
		}
	}

	return diffBuilder.String(), nil
}

func readCopilotInstructions(ctx context.Context, client *github.Client, owner, repo string) ([]string, error) {
	// Try to get the file from the repository
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, ".github/copilot-instructions.yml", nil)
	if err != nil {
		return nil, err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return nil, err
	}

	var rules ReviewRules
	if err := yaml.Unmarshal([]byte(content), &rules); err != nil {
		return nil, err
	}

	return rules.Rules, nil
}

func performCodeReview(ctx context.Context, apiKey, diff string, rules []string) (string, error) {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	rulesText := strings.Join(rules, "\n- ")

	prompt := fmt.Sprintf(`You are a code reviewer. Review the following code changes and provide feedback based EXCLUSIVELY on these rules:

- %s

Code changes:
%s

Provide a concise review with:
1. Issues found (if any) according to the rules above
2. Specific suggestions for improvement (if applicable)
3. Line references where relevant

If no issues are found, simply say "No issues found based on the review rules."`, rulesText, diff)

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: int64(4096),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return "", err
	}

	// Extract text from response
	var review strings.Builder
	for _, block := range message.Content {
		if block.Type == "text" {
			review.WriteString(block.Text)
		}
	}

	return review.String(), nil
}

func postReviewComment(ctx context.Context, client *github.Client, owner, repo string, prNum int, review string) error {
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("## 🤖 AI Code Review\n\n%s\n\n---\n*Powered by Claude via Anthropic SDK*", review)),
	}

	_, _, err := client.Issues.CreateComment(ctx, owner, repo, prNum, comment)
	return err
}
