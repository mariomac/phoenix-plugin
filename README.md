# AI Code Review GitHub Action

This GitHub Action performs automated code reviews on pull requests using Claude AI via the Anthropic SDK for Go.

## Features

- **Automatic PR Analysis**: Detects changes in pull requests and performs intelligent code review
- **Custom Review Rules**: Reads rules from `.github/copilot-instructions.md` in your repository
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

Create a `.github/copilot-instructions.md` file in your repository with your review rules. The entire file content will be passed to Claude as instructions, so you can use any format (Markdown, plain text, etc.):

```markdown
# Code Review Rules

- Check for proper error handling and ensure errors are not silently ignored
- Ensure all public functions and methods have documentation comments
- Look for potential security vulnerabilities (SQL injection, XSS, hardcoded credentials)
- Verify proper variable and function naming conventions following the project style guide
- Check for code duplication and suggest refactoring opportunities
- Ensure new features include appropriate unit tests
```

See [examples/copilot-instructions.md](examples/copilot-instructions.md) for more examples.

### 3. Add Workflow to Your Repository

Create a `.github/workflows/code-review.yml` file in your repository:

```yaml
name: AI Code Review

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  code-review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run AI Code Review
        uses: mariomac/phoenix-plugin@main
        with:
          anthropic-api-key: ${{ secrets.ANTHROPIC_API_KEY }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

**Custom rules file path:**

If you want to use a different path for your rules file:

```yaml
      - name: Run AI Code Review
        uses: mariomac/phoenix-plugin@main
        with:
          anthropic-api-key: ${{ secrets.ANTHROPIC_API_KEY }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
          rules-file: 'docs/review-rules.yml'
```

The action will automatically run on:
- New pull requests
- Updated pull requests
- Reopened pull requests

## How It Works

1. When a PR is created or updated, the action triggers
2. It fetches the PR diff to see what changed
3. It reads the complete content of your rules file (default: `.github/copilot-instructions.md`)
4. It sends the diff and the entire rules file content to Claude AI for analysis
5. Claude reviews the code based exclusively on the instructions in your rules file
6. The review is posted as a comment on the PR

**Note:** The action reads and passes the complete file content as-is to Claude, allowing you to write your review instructions in any format that works best for your team (structured Markdown, simple bullet points, detailed documentation, etc.).

## Requirements

- Anthropic API key with access to Claude models
- GitHub Actions enabled on your repository
- Docker support in GitHub Actions (enabled by default)

## Action Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `anthropic-api-key` | Anthropic API key for Claude access | Yes | - |
| `github-token` | GitHub token for API access | Yes | `${{ github.token }}` |
| `rules-file` | Path to the rules file in the repository (any text format) | No | `.github/copilot-instructions.md` |

**Note:** The entire content of the rules file is passed directly to Claude as instructions. You can use any text-based format: Markdown, plain text, or structured documentation.

## Environment Variables

The action automatically configures these environment variables:

- `ANTHROPIC_API_KEY`: Your Anthropic API key (from inputs)
- `GITHUB_TOKEN`: GitHub token for API access (from inputs)
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
