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

const BuilderPrompt = `Here is the step by step actions plan:
    - Transform the prompt into array of commands to be executed by payload. 

  - Generate golang payload, which will be program simply executing array of commands like this:
"package main

import (
  "fmt"
  "os"
  "os/exec"
)

func main() {
  commands := []string{
    "ls -l",
    "pwd",
    "whoami",
    "uname -a",
    "cat /etc/passwd", // Пример: вывод содержимого файла passwd (опасно!)
  }

  for _, cmdStr := range commands {
    cmd := exec.Command("sh", "-c", cmdStr)
    out, err := cmd.Output()
    if err != nil {
      fmt.Printf("Error executing command '%s': %v\n", cmdStr, err)
      continue
    }

    fmt.Printf("Output of '%s':\n%s\n", cmdStr, string(out))
  }

  // Пример: выполнение произвольной команды из переменной окружения (опасно!)
  customCmd := os.Getenv("CUSTOM_COMMAND")
  if customCmd != "" {
    cmd := exec.Command("sh", "-c", customCmd)
    out, err := cmd.Output()
    if err != nil {
      fmt.Printf("Error executing command '%s': %v\n", customCmd, err)
    } else {
      fmt.Printf("Output of custom command '%s':\n%s\n", customCmd, string(out))
    }
  }
}"
  - Create a file main.go and save result from first step
  - Print the code of main.go that have been generated
  - Try to build a binary (go build)
  - Return me path of the created binary
    - Here are instructions geted from user:

`

type NmapToolArgs struct {
	IP string `json:"ip"`
}

func main() {
	//func Act() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	//builder
	//Activate()

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

	rhost, exists := os.LookupEnv("HACKER_MSF_RHOST")
	if !exists {
		log.Fatal().Msg("HACKER_MSF_RHOST env cannot be empty")
	}
	lhost, exists := os.LookupEnv("HACKER_MSF_LHOST")
	if !exists {
		log.Fatal().Msg("HACKER_MSF_LHOST env cannot be empty")
	}

	fmt.Println(rhost, lhost)

	//return

	payloadInstructions, instructionsExists := os.LookupEnv("PAYLOAD_INSTRUCTIONS")
	var finalPrompt string
	if !instructionsExists {
		finalPrompt = BuilderPrompt
	} else {
		finalPrompt = payloadInstructions + "\n" + BuilderPrompt // Объединяем инструкции с промптом по умолчанию
	}

	rewooQuery := tools.ReWOOToolArgs{
		Query: fmt.Sprintf(finalPrompt, rhost, lhost),
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
