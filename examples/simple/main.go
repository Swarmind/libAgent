package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Swarmind/libagent/pkg/agent/simple"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms/openai"
)

/*
	This example shows how to initialize and use a simple agent without any tools or specific stuff.
*/

const Prompt = `This is a test. Write OK in response.`

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()
	agent := simple.Agent{}

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

	result, err := agent.SimpleRun(ctx,
		Prompt, config.ConifgToCallOptions(cfg.DefaultCallOptions)...,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("agent run")
	}
	fmt.Println(util.RemoveThinkTag(result))
}
