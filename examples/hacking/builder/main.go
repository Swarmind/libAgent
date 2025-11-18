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
  This example shows usage of command executor to create a binary with generated code.
*/

const BuilderPrompt = `Here is the step by step actions plan, use the command executor tool to create the binary at the end:
	- Generate a simple go hello world code
	- Write it to the file (use echo -e or printf with single quotes since echo in bash did not interpret escaped symbols by default, and there is multiple symbols needed to be escaped using single quotes as well)
	- Initialize go module
	- Build it
	- Use pwd to get the working directory
	- Return me the path to built binary
`

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig("")
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

	// Do not cleanup to use the final result
	// defer func() {
	// 	if err := toolsExecutor.Cleanup(); err != nil {
	// 		log.Fatal().Err(err).Msg("tools executor cleanup")
	// 	}
	// }()

	rewooQuery := tools.ReWOOToolArgs{
		Query: BuilderPrompt,
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
