package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/tmc/langchaingo/llms"
)

var MsfSearchToolDefinition = llms.FunctionDefinition{ // Переименовано на MsfSearchToolDefinition
	Name:        "msf_search", // Переименовано на msf_search
	Description: "Executes Metasploit search queries provided in a list.",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"queries": map[string]any{
				"type":        "array",
				"description": "A list of Metasploit search queries to execute.",
			},
		},
	},
}

type MsfSearchTool struct{} // Переименовано на MsfSearchTool

type MsfSearchToolArgs struct { // Переименовано на MsfSearchToolArgs
	Queries []string `json:"queries"`
}

// Call executes the Metasploit search commands with the given queries.
func (s MsfSearchTool) Call(ctx context.Context, input string) (string, error) { // Переименовано на MsfSearchTool
	msfToolArgs := MsfSearchToolArgs{} // Переименовано на MsfSearchToolArgs
	if err := json.Unmarshal([]byte(input), &msfToolArgs); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var results []string
	for _, query := range msfToolArgs.Queries {
		cmd := exec.Command("msfconsole", "-q", "-x", "search "+query+"; exit")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to execute msfconsole for query '%s': %w", query, err)
		}

		results = append(results, string(output))
	}

	// Формирование ответа для LangChainGo
	response := fmt.Sprintf("Metasploit search results: %v", results)
	return response, nil
}

func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			if cfg.MsfDisable {
				return nil, nil
			}

			return &tools.ToolData{
				Definition: MsfSearchToolDefinition,
				Call:       MsfSearchTool{}.Call,
			}, nil
		},
	)
}
