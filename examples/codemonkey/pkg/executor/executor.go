package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CliGenerator(task string) string {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("new config")
	}

	ctx := context.Background()

	toolsToWhitelist := []string{
		tools.ReWOOToolDefinition.Name,
		tools.CommandExecutorDefinition.Name,
	}

	toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))
	if err != nil {
		log.Fatal().Err(err).Msg("new tools executor")
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			log.Fatal().Err(err).Msg("tools executor cleanup")
		}
	}()

	rewooQuery := tools.ReWOOToolArgs{
		Query: CreatePrompt(task),
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

	if result == "" {
		log.Fatal().Msg("main empty result")
	}

	return result

}

func CreatePrompt(task string) string {
	prompt := `You are an AI command execution assistant. Your purpose is to generate CLI commands to accomplish tasks.

	RULES:
	1. Generate ONLY valid Unix/Linux CLI commands
	2. Output MUST be plain commands without any explanations
	3. Separate multiple commands with newlines (\n)
	4. Commands must be escaped properly for shell execution
	5. If task requires sequential operations, use '&&' or proper piping
	6. Never generate interactive or destructive commands
	
	EXAMPLE:
	For task "list files": ls -la
	For task "create file and show content": touch file.txt && cat file.txt
	
	TASK: ` + task + `
	
	GENERATED COMMANDS:`

	return prompt
}

func ExecuteCommands(commandStr string) error {
	commands := strings.Split(commandStr, "\n")
	var nonEmptyCommands []string
	for _, cmd := range commands {
		if trimmed := strings.TrimSpace(cmd); trimmed != "" {
			nonEmptyCommands = append(nonEmptyCommands, trimmed)
		}
	}

	if len(nonEmptyCommands) == 0 {
		return nil
	}

	script := strings.Join(nonEmptyCommands, "\n")
	cmd := exec.Command("sh", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script failed: %s\nError: %s\nOutput: %s", script, err, out)
	}

	fmt.Printf("Script Output:\n%s\n", out)
	return nil
}
