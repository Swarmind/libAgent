package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
	This example shows usage of command executor with rewoo tool, which are whitelisted.
*/

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
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(
		tools.ReWOOToolDefinition.Name,
		tools.CommandExecutorDefinition.Name,
	))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

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
		log.Fatal().Err(err).Msg("rewoo tool call")
	}

	fmt.Println(result)
}
