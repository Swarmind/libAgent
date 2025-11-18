package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Prompt = `Please lookup for address of "%s". Use different tools ` +
	`such as dig or nmap (don't forget -F flag, since execution timeout is 30 seconds) through command executor, ` +
	`https://internetdb.shodan.io/<ip> utility via webReader tool, ` +
	`shodan dedicated tool to find out open ports and additional addresses of the target. ` +
	`After that create metasploit queries for each of found open ports if any using msf_search tool.`

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
		tools.WebReaderDefinition.Name,
		tools.ShodanToolDefinition.Name,
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

	addr, ok := os.LookupEnv("IPLOOKUP_ADDR")
	if !ok {
		fmt.Println("IPLOOKUP_ADDR env is empty.\nEnter target ADDR (IP/CIDR/URL):")
		_, err = fmt.Scanln(&addr)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read stdin for target ADDR")
		}
	}
	addr = strings.TrimSpace(addr)
	validIP := net.ParseIP(addr) != nil
	validCIDR := true
	_, _, cidrErr := net.ParseCIDR(addr)
	if cidrErr != nil {
		validCIDR = false
	}
	validURL := true
	_, urlErr := url.Parse(addr)
	if urlErr != nil {
		validURL = false
	}

	if !validIP && !validCIDR && !validURL {
		log.Fatal().Msg("validte provided address")
	}

	fmt.Printf(
		"Provided addr: %s, ip: %t, ip/cidr: %t, url: %t\n",
		addr,
		validIP,
		validCIDR,
		validURL,
	)

	rewooQuery := tools.ReWOOToolArgs{
		Query: fmt.Sprintf(Prompt, addr),
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
