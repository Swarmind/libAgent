package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
	This example shows how to do tool call multiple times for different config values.
*/

var ModelList = []string{
	"qwen3-32b",
	"qwen3-30b-a3b",
	"qwen3-30b-a3b-abliterated",
	"qwen3-14b",
	"qwen3-8b",
	"mlabonne_qwen3-8b-abliterated",
	"josiefied-qwen3-8b-abliterated-v1",
	"big-tiger-gemma-27b-v1",
	"gemma-3-27b-it",
	"gemma-3-12b-it",
	"tiger-gemma-9b-v1-i1",
	"wizard-vicuna-13b-uncensored",
	"deepseek-r1-distill-llama-8b",
	"deepseek-r1-distill-qwen-1.5b",
	"deepseek-r1-distill-qwen-7b",
	"deepseek-r1-distill-qwen-14b",
	"deepseek-r1-distill-qwen-32b",
	"deepseek-r1-distill-llama-70b",
	"huihui-ai_deepseek-r1-distill-llama-70b-abliterated",
}

const Prompt = `Here is the step by step actions plan:
	- Create a file with 'banana' content
	- Read from it and verify it is indeed contains 'banana'
	- Modify it to 'bannana'
	- Read from it again and verify it is indeed containes 'bannana'
	- Check if 172.86.66.189 host has open 443 and 13370 ports
	- Get weather from wttr.in service with format=3 option
	- Write a hello world python script
	- Execute it and verify it is working
	- Create a directory "testrepo" and move that script inside
	- Initialize a git repository in that directory
	- Configure git locally to use your (LLM) arbitrary name and email for the commit
	- Do a commit

As a result of future plan solving - write me a report as a list like this:
	- create file with contents 'banana' [echo 'banana' > file.txt]: OK
	- verify file contains 'banana' [cat file.txt]: banana
	- ...
`

func main() {
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	for idx, model := range ModelList {
		log.Info().Msgf("Using model: %s", model)
		cfg.Model = model

		toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(
			tools.ReWOOToolDefinition.Name,
			tools.CommandExecutorDefinition.Name,
		))
		if err != nil {
			log.Fatal().Err(err).Msg("new tools executor")
		}

		rewooQuery := tools.ReWOOToolArgs{
			Query: Prompt,
		}
		rewooQueryBytes, err := json.Marshal(rewooQuery)
		if err != nil {
			log.Fatal().Err(err).Msg("json marhsal rewooQuery")
		}

		result, err := toolsExecutor.CallTool(ctx,
			tools.ReWOOToolDefinition.Name,
			string(rewooQueryBytes),
		)
		if err != nil {
			log.Warn().Err(err).Msg("rewoo tool call")
			if idx != len(ModelList)-1 {
				log.Info().Msgf("Sleeping 2 minutes so LocalAI watchdog unloads model")
				time.Sleep(time.Minute * 2)
			}
			continue
		}

		fmt.Println(result)

		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}

		if idx != len(ModelList)-1 {
			log.Info().Msgf("Sleeping 2 minutes so LocalAI watchdog unloads model")
			time.Sleep(time.Minute * 2)
		}
	}
}
