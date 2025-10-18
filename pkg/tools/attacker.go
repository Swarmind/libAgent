package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	//"strings"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/tmc/langchaingo/llms"
)

var ExploitToolDefinition = llms.FunctionDefinition{
	Name:        "exploit",
	Description: "Executes Metasploit exploit against a target using provided module and options.",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"module": map[string]any{
				"type":        "string",
				"description": "Metasploit exploit module to use (e.g., 'exploit/multi/handler').",
			},
			"options": map[string]any{
				"type":        "object",
				"description": "Key-value pairs for Metasploit options (e.g., {'RHOSTS': '192.168.1.10', 'LHOST': '192.168.1.5'}).",
			},
		},
	},
}

type ExploitTool struct{}

type ExploitToolArgs struct {
	Module  string            `json:"module"`
	Options map[string]string `json:"options"`
}

func (s *ExploitTool) Call(ctx context.Context, input string) (string, error) {
	exploitArgs := ExploitToolArgs{}
	if err := json.Unmarshal([]byte(input), &exploitArgs); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	cmd := []string{"msfconsole"}
	for key, value := range exploitArgs.Options {
		cmd = append(cmd, "-o", key+"="+value)
	}
	cmd = append(cmd, "use", exploitArgs.Module)
	cmd = append(cmd, "run")

	output, err := exec.Command("msfconsole", cmd...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute msfconsole: %w - output: %s", err, string(output))
	}

	return string(output), nil
}

func init() {
	globalToolsRegistry = append(globalToolsRegistry, func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
		if cfg.ExploitDisable {
			return nil, nil
		}

		tool := ExploitTool{}

		return &tools.ToolData{
			Definition: ExploitToolDefinition,
			Call:       tool.Call,
		}, nil
	})
}
