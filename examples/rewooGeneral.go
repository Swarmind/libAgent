package main

import (
	"context"
	"libagent/pkg/agent/rewoo"
	"libagent/pkg/config"
	"libagent/pkg/tools"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms/openai"
)

const GenericTestPrompt = `Using semantic search tool, which can search across various code from the project
collections find out the telegram library name in the code file contents for the project called "Hellper".
Extract it from the given code and use a web search to find the pkg.go.dev documentation for it.
Give me the URL for it.`

const InnerReWOOTestPrompt = `This is a test of your ability to use various tools.
I need you to simply call the tool called "rewoo" for creating subtasks with the next prompt and return the result.
Prompt:
` + GenericTestPrompt

const HackerTestPrompt = `Using duck-duck go search tool, which can search across web, create a cyberattack vector. 
First analyse user mission and point a target. 
Then try to get some information about target by web.
Then create cyberattack vector to this target. 
Then think what kind of CVE's can be used to actual perform attack on the target
Then use search and find any relevant information about this CVE in a web
Return user this plan.

User Mission:
`

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("main config.NewConfig")
	}
	if cfg.AIURL == "" {
		log.Fatal().Msg("main empty OpenAI URL")
	}
	if cfg.ApiToken == "" {
		log.Fatal().Msg("main empty OpenAI Token")
	}
	if cfg.Model == "" {
		log.Fatal().Msg("main empty model")
	}
	rewooAgent := rewoo.Agent{}
	ctx := context.WithValue(context.Background(), "ReWOOAgent", &rewooAgent)

	llm, err := openai.New(
		openai.WithBaseURL(cfg.AIURL),
		openai.WithToken(cfg.ApiToken),
		openai.WithModel(cfg.Model),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("main openai.New")
	}
	rewooAgent.LLM = llm

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("main tools.NewToolsExecutor")
	}
	rewooAgent.ToolsExecutor = toolsExecutor

    result, err := rewooAgent.SimpleRun(ctx, HackerTestPrompt)
    if err != nil {
        log.Fatal().Err(err).Msg("main rewooAgent.SimpleRun")
    }

    if result == "" {
        log.Fatal().Msg("main empty result")
    }

    log.Info().Str("result", result).Msg("examples/hacker")
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Msgf("%s env is empty", key)
	}

	return val
}
