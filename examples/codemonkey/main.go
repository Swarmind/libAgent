package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultRepoURL = "https://github.com/JackBekket/Reflexia"

/*
	This example uses the ReWOO agent to manage a GitHub issue end-to-end,
	inspired by the prompting style of examples/generic/main.go.
	The prompt guides ReWOO on tool usage for various sub-tasks.

	Default repository: https://github.com/JackBekket/Reflexia
	Default semantic search collection (if default repo is used): Reflexia
*/

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	repoURL := flag.String("repo", defaultRepoURL, "Repository URL")
	issueInput := flag.String("issue", "", "Issue number, full URL, or descriptive text (required)")
	additionalCtx := flag.String("additionalContext", "", "Optional additional context for the agent.")
	flag.Parse()

	if *issueInput == "" {
		log.Fatal().Msg("Issue input is required. Use -issue flag.")
		return
	}
	log.Debug().Str("issue_input", *issueInput).Msg("Received issue input")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application config.")
		return
	}

	ctx := context.Background()

	toolsToWhitelist := []string{
		tools.ReWOOToolDefinition.Name,
		tools.CommandExecutorDefinition.Name,
		tools.SemanticSearchDefinition.Name,
		tools.WebReaderDefinition.Name,
	}
	log.Debug().Interface("whitelisted_tools", toolsToWhitelist).Msg("Initializing ToolsExecutor.")

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create tools executor.")
		return
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Error().Err(err).Msg("Error during tools executor cleanup.")
		}
	}()

	var repoName string
	var ownerAndRepo string
	trimmedRepoURL := strings.TrimSuffix(*repoURL, ".git")

	if trimmedRepoURL == strings.TrimSuffix(defaultRepoURL, ".git") {
		repoName = "Reflexia"
		ownerAndRepo = "JackBekket/Reflexia"
	} else {
		repoName = filepath.Base(trimmedRepoURL)
		parts := strings.Split(trimmedRepoURL, "/")
		if len(parts) >= 2 {
			ownerAndRepo = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
		} else {
			ownerAndRepo = repoName
		}
	}
	log.Debug().Str("derived_repo_name", repoName).Str("derived_owner_repo", ownerAndRepo).Msg("Repository details derived.")

	additionalContextString := ""
	if *additionalCtx != "" {
		additionalContextString = fmt.Sprintf("\n\nAdditional Context Provided by User:\n%s", *additionalCtx)
	}

	taskDescription := fmt.Sprintf(`
You are an AI Code Monkey tasked with resolving GitHub issue "%s" for repository %s. Your goal is to understand the issue and implement the fix.
%s

Follow these steps:

1. Understand the Issue:
   - If the issue is a URL, use WebReaderTool to fetch its content (you are allowed to access links)
   - Use LLM to analyze the issue content and extract key requirements
   - Use SemanticSearchTool to find relevant files in the repository

2. Set Up Development Environment:
   - Use CommandExecutor to run: git clone %s
   - Use CommandExecutor to change to repo directory: cd %s
   - Use CommandExecutor to mark the directory as safe: cd %s && git config --global --add safe.directory "$(pwd)"
   - Use CommandExecutor to view file structure: tree or ls -R

3. Find and Modify Files:
   - Use SemanticSearchTool to find relevant files and their contents
   - For each file that needs changes:
     * First, view the current file content: cd %s && cat <file>
     * Create a backup: cd %s && cp <file> <file>.bak
     * Make the required changes using CommandExecutor:
       - For small changes: cd %s && echo 'new content' > <file>
       - For larger changes: cd %s && cat > <file> << 'EOL'
         new content
         multiple lines
         EOL
     * Verify changes: cd %s && cat <file>
     * Verify git status: cd %s && git status
     * If changes are not showing in git status, make sure you're modifying the original file, not the backup

4. After Making Changes:
   - Verify all changes are present: cd %s && git status
   - If changes are not showing, check:
     * Are you in the correct directory?
     * Are you modifying the original files (not .bak files)?
     * Are the changes being written to the correct paths?
   - If needed, use git diff to see what changes were made: cd %s && git diff

Note: 
- Use CommandExecutor to run all commands, including viewing and modifying files
- Use tree or ls to explore the repository structure
- Use cat to view file contents before and after changes
- For multi-line changes, use heredoc syntax (<< 'EOL')
- Always verify changes after making them
- Make sure to modify the original files, not just create backups
- If changes aren't showing in git status, double-check your file paths and commands
`, *issueInput, *repoURL, additionalContextString, *repoURL, repoName, repoName, repoName, repoName, repoName, repoName, repoName, repoName, repoName, repoName)

	log.Info().Str("repository", *repoURL).Str("issue", *issueInput).Msg("Initiating ReWOO agent for Code Monkey task.")
	if *additionalCtx != "" {
		log.Info().Str("additional_context_provided_length", fmt.Sprintf("%d chars", len(*additionalCtx))).Msg("Additional context will be used.")
	}
	currentTaskDescLength := len(taskDescription)
	log.Debug().Str("rewoo_task_description_length", fmt.Sprintf("%d chars", currentTaskDescLength)).Msg("Constructed ReWOO Task for Code Monkey")

	if currentTaskDescLength > 4000 {
		log.Warn().Msgf("ReWOO task description is %d chars, which is over the 4000 char suggestion. This might lead to issues.", currentTaskDescLength)
	}

	rewooQueryArgs := tools.ReWOOToolArgs{
		Query: taskDescription,
	}
	rewooQueryBytes, err := json.Marshal(rewooQueryArgs)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal ReWOO query arguments.")
		return
	}

	result, err := toolsExecutor.CallTool(ctx,
		tools.ReWOOToolDefinition.Name,
		string(rewooQueryBytes),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("ReWOO agent task failed.")
		return
	}

	fmt.Println("\n--- ReWOO Agent Task Execution Report ---")
	fmt.Println(result)
	fmt.Println("\n--- End of ReWOO Agent Task ---")
	fmt.Println("\nOperational Reminders:")
	fmt.Println("1. Ensure 'git' and GitHub CLI ('gh') are installed, configured, and available in the system PATH.")
	fmt.Println("2. 'gh' CLI must be authenticated with GitHub (e.g., via 'gh auth login') for PR creation.")
	fmt.Println("3. For 'SemanticSearchTool', a 'pgvector' database with a collection named after the repository (e.g., 'Reflexia') must be populated with relevant code embeddings.")
}
