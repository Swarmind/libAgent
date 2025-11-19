package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Swarmind/libagent/pkg/agent/generic"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms/openai"
)

/*
	This example shows how to initialize and use a generic tools enabled agent.
*/

const Prompt = `Please use rewoo tool with the next prompt:
Using semantic search tool, which can search across various code from the project
collections find out the telegram library name in the code file contents for the project called "Hellper".
Extract it from the given code and use a web search to find the pkg.go.dev documentation for it.
Give me the URL for it.`

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()
	agent := generic.Agent{}

	llm, err := openai.New(
		openai.WithBaseURL(cfg.AIURL),
		openai.WithToken(cfg.AIToken),
		openai.WithModel(cfg.Model),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("new openai api llm")
	}
	agent.LLM = llm

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(
		tools.ReWOOToolDefinition.Name,
		tools.SemanticSearchDefinition.Name,
		tools.DDGSearchDefinition.Name,
		tools.WebReaderDefinition.Name,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	agent.ToolsExecutor = toolsExecutor
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

	result, err := agent.SimpleRun(ctx,
		Prompt, config.ConifgToCallOptions(cfg.DefaultCallOptions)...,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("agent run")
	}
	fmt.Println(result)
}
