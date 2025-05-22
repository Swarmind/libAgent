package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Swarmind/libagent/pkg/agent/codemonkey" // Import the CodeMonkeyAgent package
	"github.com/Swarmind/libagent/pkg/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
	This example demonstrates how to use the Code Monkey Agent to automatically
	resolve issues in a codebase and create pull requests with the fixes.

	Usage:
		go run main.go -repo="https://github.com/owner/repo.git" -issue="123"
		or
		go run main.go -repo="https://github.com/StarumDDD/Personalized-portfolio-cite-TS" -issue="https://github.com/StarumDDD/Personalized-portfolio-cite-TS/issues/1"
*/

func main() {
	// Set up logging
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Parse command line arguments
	repoURL := flag.String("repo", "", "Repository URL (required)")
	issueInput := flag.String("issue", "", "Issue number or URL (required)")
	flag.Parse()

	if *repoURL == "" || *issueInput == "" {
		log.Fatal().Msg("Repository URL and issue are required. Use -repo and -issue flags")
		return // Ensure we exit if fatal error occurs
	}

	// Initialize config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new config")
		return
	}

	// Create context
	ctx := context.Background()

	// Create a new CodeMonkeyAgent
	agent, err := codemonkey.NewCodeMonkeyAgent(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create CodeMonkeyAgent")
		return
	}
	defer func() {
		if err := agent.Cleanup(); err != nil {
			log.Error().Err(err).Msg("Failed during agent cleanup")
		}
	}()

	log.Info().Msgf("Attempting to resolve issue '%s' in repository '%s'", *issueInput, *repoURL)

	// Call the ResolveIssue method of the CodeMonkeyAgent
	result, err := agent.ResolveIssue(ctx, *repoURL, *issueInput)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to resolve issue")
		return
	}

	fmt.Println("\nCode Monkey Agent Task Completed:")
	fmt.Println(result)
	fmt.Println("\nDebug Info from original example (for user reference):")
	fmt.Println("1. Make sure you have git installed and configured (including GitHub CLI 'gh' for PR creation).")
	fmt.Println("2. The repository will be cloned in a temporary directory (managed by the agent).")
	fmt.Println("3. All git operations will be performed in the cloned repository.")
}
