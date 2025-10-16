# AI Code Review GitHub Action

This GitHub Action performs automated code reviews on pull requests using Claude AI via the Anthropic SDK for Go.

## Features

- **Automatic PR Analysis**: Detects changes in pull requests and performs intelligent code review
- **Custom Review Rules**: Reads rules from `.github/copilot-instructions.yml` in your repository
- **AI-Powered Suggestions**: Uses Claude to provide specific code improvement suggestions
- **Inline Feedback**: Posts review comments directly on the pull request

## Setup

### 1. Add Anthropic API Key

Add your Anthropic API key as a repository secret:

1. Go to your repository Settings > Secrets and variables > Actions
2. Click "New repository secret"
3. Name: `ANTHROPIC_API_KEY`
4. Value: Your Anthropic API key

### 2. Create Review Rules

Create a `.github/copilot-instructions.yml` file in your repository with your review rules:

```yaml
rules:
  - Check for proper error handling
  - Ensure all functions have documentation comments
  - Look for security vulnerabilities
  - Verify proper variable naming conventions
  - Check for code duplication
  - Ensure tests are included for new features
```

### 3. Copy the Workflow

The GitHub Action workflow is already set up in `.github/workflows/code-review.yml`. It will automatically run on:
- New pull requests
- Updated pull requests
- Reopened pull requests

## How It Works

1. When a PR is created or updated, the action triggers
2. It fetches the PR diff to see what changed
3. It reads your custom review rules from `.github/copilot-instructions.yml`
4. It sends the diff and rules to Claude AI for analysis
5. Claude reviews the code based exclusively on your specified rules
6. The review is posted as a comment on the PR

## Requirements

- Anthropic API key with access to Claude models
- GitHub Actions enabled on your repository
- Go 1.21 or later (handled automatically by the action)

## Environment Variables

The action uses these environment variables (automatically set by GitHub Actions):

- `ANTHROPIC_API_KEY`: Your Anthropic API key (secret)
- `GITHUB_TOKEN`: GitHub token for API access (automatic)
- `GITHUB_REPOSITORY`: Repository name (automatic)
- `PR_NUMBER`: Pull request number (automatic)

## Docker Support

You can also run the code reviewer using Docker:

```bash
docker build -t code-reviewer .
docker run \
  -e ANTHROPIC_API_KEY=your-key \
  -e GITHUB_TOKEN=your-token \
  -e GITHUB_REPOSITORY=owner/repo \
  -e PR_NUMBER=123 \
  code-reviewer
```

## Local Development

```bash
# Install dependencies
go mod download

# Set environment variables
export ANTHROPIC_API_KEY=your-key
export GITHUB_TOKEN=your-token
export GITHUB_REPOSITORY=owner/repo
export PR_NUMBER=123

# Run
go run main.go
```
