package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
	This example shows usage of command executor with rewoo tool, which are whitelisted.
*/

const Prompt = `Please scan 192.168.1.0 and 192.168.1.1  for open ports and generate Metasploit search queries for any found services.`

type NmapToolArgs struct {
	IP string `json:"ip"`
}

// struct для описания найденного порта
type PortInfo struct {
	Port    string
	State   string
	Service string
}

func parseNmapPorts(nmapOutput string) []PortInfo {
	var ports []PortInfo
	re := regexp.MustCompile(`(\d+)\/tcp\s+(\w+)\s+(.+?)\s*`)
	matches := re.FindAllStringSubmatch(nmapOutput, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			ports = append(ports, PortInfo{
				Port:    match[1],
				State:   strings.TrimSpace(match[2]),
				Service: strings.TrimSpace(match[3]),
			})
		}
	}

	return ports
}

func generateMsfQueries(ports []PortInfo) []string {
	var queries []string
	for _, port := range ports {
		if strings.ToLower(port.State) == "open" {
			queries = append(queries, fmt.Sprintf("type:exploit name:%s", port.Service))
			queries = append(queries, fmt.Sprintf("port %s", port.Port))
		}
	}

	return queries
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
