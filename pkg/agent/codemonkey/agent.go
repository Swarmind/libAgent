package codemonkey

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	itools "github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	ptools "github.com/Swarmind/libagent/pkg/tools"

	"github.com/tmc/langchaingo/llms/openai"
	"github.com/rs/zerolog"
)

var log = zerolog.New(os.Stderr).With().Timestamp().Logger()

// CodeMonkeyAgent implements the Code Monkey flow for automated issue resolution
type CodeMonkeyAgent struct {
	LLM           *openai.LLM
	ToolsExecutor *itools.ToolsExecutor
	WorkDir       string
}

// NewCodeMonkeyAgent creates a new Code Monkey Agent instance
func NewCodeMonkeyAgent(ctx context.Context, cfg config.Config) (*CodeMonkeyAgent, error) {
	llm, err := openai.New(
		openai.WithBaseURL(cfg.AIURL),
		openai.WithToken(cfg.AIToken),
		openai.WithModel(cfg.Model),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}

	toolsExecutor, err := ptools.NewToolsExecutor(ctx, cfg, ptools.WithToolsWhitelist(
		ptools.ReWOOToolDefinition.Name,
		ptools.SemanticSearchDefinition.Name,
		ptools.CommandExecutorDefinition.Name,
		ptools.DDGSearchDefinition.Name,
		ptools.WebReaderDefinition.Name,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create tools executor: %w", err)
	}

	workDir, err := os.MkdirTemp("", "codemonkey_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	return &CodeMonkeyAgent{
		LLM:           llm,
		ToolsExecutor: toolsExecutor,
		WorkDir:       workDir,
	}, nil
}

// executeGitCommand executes a git command using the command executor tool
func (a *CodeMonkeyAgent) executeGitCommand(ctx context.Context, repoPath string, command string) error {
	cmd := ptools.CommandExecutorArgs{
		Command: fmt.Sprintf("cd %s && %s", repoPath, command),
	}
	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	_, err = a.ToolsExecutor.CallTool(ctx,
		ptools.CommandExecutorDefinition.Name,
		string(cmdBytes),
	)
	return err
}

// ResolveIssue implements the Code Monkey flow for issue resolution
func (a *CodeMonkeyAgent) ResolveIssue(ctx context.Context, repoURL string, issueInput string) (string, error) {
	// Step 1: Get issue text from GitHub
	issueText, err := a.getIssueText(ctx, repoURL, issueInput)
	if err != nil {
		return "", fmt.Errorf("failed to get issue text: %w", err)
	}

	// Extract owner and repo from URL for PR creation
	parts := strings.Split(strings.TrimSuffix(repoURL, ".git"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid repository URL: %s", repoURL)
	}
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Step 2: Clone the repository
	repoName := filepath.Base(strings.TrimSuffix(repoURL, ".git"))
	repoPath := filepath.Join(a.WorkDir, repoName)

	if err := a.executeGitCommand(ctx, a.WorkDir, fmt.Sprintf("git clone %s %s", repoURL, repoName)); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Step 3: Create a new branch
	// Use the first line of the issue text for the branch name
	firstLine := strings.Split(issueText, "\n")[0]
	// Clean the title for branch name
	cleanTitle := strings.ToLower(firstLine)
	// Remove special characters and replace spaces with hyphens
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	cleanTitle = reg.ReplaceAllString(cleanTitle, "")
	cleanTitle = strings.TrimSpace(cleanTitle)
	cleanTitle = strings.ReplaceAll(cleanTitle, " ", "-")
	
	// Take up to 30 characters, but don't exceed the length of the title
	titleLength := len(cleanTitle)
	if titleLength > 30 {
		titleLength = 30
	}
	branchName := fmt.Sprintf("fix/%s", cleanTitle[:titleLength])

	// Log the branch name for debugging
	log.Debug().Str("branchName", branchName).Msg("Creating branch")

	// Ensure we're in the repository directory and on main branch
	if err := a.executeGitCommand(ctx, repoPath, "git checkout main"); err != nil {
		return "", fmt.Errorf("failed to checkout main branch: %w", err)
	}

	// Create and checkout the new branch
	branchCmd := fmt.Sprintf("git checkout -b %s", branchName)
	log.Debug().Str("command", branchCmd).Msg("Executing branch creation command")
	if err := a.executeGitCommand(ctx, repoPath, branchCmd); err != nil {
		return "", fmt.Errorf("failed to create branch: %w", err)
	}

	// Step 4: Context Retrieval using semantic search
	semanticSearchArgs := ptools.SemanticSearchArgs{
		Query:      issueText,
		Collection: "codebase",
	}
	semanticSearchBytes, err := json.Marshal(semanticSearchArgs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal semantic search args: %w", err)
	}

	context, err := a.ToolsExecutor.CallTool(ctx,
		ptools.SemanticSearchDefinition.Name,
		string(semanticSearchBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve context: %w", err)
	}

	// Step 5: Solution Generation using ReWOO
	rewooQuery := ptools.ReWOOToolArgs{
		Query: fmt.Sprintf(`Given the following issue and context, generate a solution plan:
Issue: %s
Context: %s

Generate a detailed plan to resolve this issue. The plan should include:
1. Analysis of the problem
2. Required changes
3. Implementation steps
4. Testing strategy
5. Validation steps`, issueText, context),
	}
	rewooQueryBytes, err := json.Marshal(rewooQuery)
	if err != nil {
		return "", fmt.Errorf("failed to marshal rewoo query: %w", err)
	}

	solution, err := a.ToolsExecutor.CallTool(ctx,
		ptools.ReWOOToolDefinition.Name,
		string(rewooQueryBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate solution: %w", err)
	}

	// Step 6: Plan Execution
	executionQuery := ptools.ReWOOToolArgs{
		Query: fmt.Sprintf(`Execute the following solution plan in the repository at %s:
%s

Use the command executor tool to implement the changes.
Verify each step and report progress.`, repoPath, solution),
	}
	executionQueryBytes, err := json.Marshal(executionQuery)
	if err != nil {
		return "", fmt.Errorf("failed to marshal execution query: %w", err)
	}

	result, err := a.ToolsExecutor.CallTool(ctx,
		ptools.ReWOOToolDefinition.Name,
		string(executionQueryBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to execute plan: %w", err)
	}

	// Step 7: Commit and push changes
	commitMessage := fmt.Sprintf("Fix: %s", firstLine)
	if err := a.executeGitCommand(ctx, repoPath, fmt.Sprintf("git add . && git commit -m \"%s\"", commitMessage)); err != nil {
		return "", fmt.Errorf("failed to commit changes: %w", err)
	}

	if err := a.executeGitCommand(ctx, repoPath, fmt.Sprintf("git push origin %s", branchName)); err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}

	// Step 8: Create pull request
	prBody := fmt.Sprintf(`This PR fixes the issue: %s

Implementation details:
%s

Changes made:
%s`, issueInput, solution, result)

	// Create PR using GitHub CLI
	prCmd := fmt.Sprintf("gh pr create --repo %s/%s --title \"%s\" --body \"%s\" --base main --head %s",
		owner, repo, commitMessage, prBody, branchName)
	
	if err := a.executeGitCommand(ctx, repoPath, prCmd); err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	return fmt.Sprintf("Pull request has been created for branch %s.\n\nImplementation details:\n%s", 
		branchName, result), nil
}

// getIssueText fetches the issue text from GitHub
func (a *CodeMonkeyAgent) getIssueText(ctx context.Context, repoURL string, issueInput string) (string, error) {
	// Extract owner and repo from URL
	parts := strings.Split(strings.TrimSuffix(repoURL, ".git"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid repository URL: %s", repoURL)
	}
	owner := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Extract issue number
	var issueNumber string
	if strings.Contains(issueInput, "/issues/") {
		// Handle full issue URL
		parts := strings.Split(issueInput, "/issues/")
		issueNumber = parts[len(parts)-1]
	} else {
		// Handle just the issue number
		issueNumber = issueInput
	}

	// Use git command to get issue text
	cmd := ptools.CommandExecutorArgs{
		Command: fmt.Sprintf("gh issue view %s --repo %s/%s --json title,body --jq '.title + \"\\n\\n\" + .body'", issueNumber, owner, repo),
	}
	cmdBytes, err := json.Marshal(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to marshal command: %w", err)
	}

	issueText, err := a.ToolsExecutor.CallTool(ctx,
		ptools.CommandExecutorDefinition.Name,
		string(cmdBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get issue text: %w", err)
	}

	// Log the issue text for debugging
	log.Debug().Str("issueText", issueText).Msg("Retrieved issue text")

	return issueText, nil
}

// Cleanup performs necessary cleanup operations
func (a *CodeMonkeyAgent) Cleanup() error {
	if a.ToolsExecutor != nil {
		if err := a.ToolsExecutor.Cleanup(); err != nil {
			return err
		}
	}
	if a.WorkDir != "" {
		return os.RemoveAll(a.WorkDir)
	}
	return nil
}
