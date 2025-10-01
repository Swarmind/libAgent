package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/tmc/langchaingo/llms"
)

var MsfSearchToolDefinition = llms.FunctionDefinition{
	Name:        "msf_search",
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

type MsfSearchTool struct {
	executable   string
	argsTemplate string // template with one %s for the query -- `search %s; exit`
}

type MsfSearchToolArgs struct {
	Queries []string `json:"queries"`
}

var (
	msfExecutable   = "msfconsole"
	msfArgsTemplate = "search %s; exit" // will be passed to msfconsole like "-q -x"
)

func SetMsfCommand(executable string, argsTemplate string) {
	if executable != "" {
		msfExecutable = executable
	}
	if argsTemplate != "" {
		msfArgsTemplate = argsTemplate
	}
}

// constructs args as []string{"-q", "-x", fmt.Sprintf(msfArgsTemplate, query)}
func (s MsfSearchTool) Call(ctx context.Context, input string) (string, error) {
	msfToolArgs := MsfSearchToolArgs{}
	if err := json.Unmarshal([]byte(input), &msfToolArgs); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var results []map[string]string
	for _, query := range msfToolArgs.Queries {
		cmdArg := fmt.Sprintf(msfArgsTemplate, query)
		args := []string{"-q", "-x", cmdArg}

		execName := s.executable
		if execName == "" {
			execName = msfExecutable
		}

		cmd := exec.Command(execName, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to execute %s %v: %w\noutput: %s", execName, args, err, string(output))
		}

		results = append(results, map[string]string{
			"query":  query,
			"output": string(output),
		})
	}

	respBytes, err := json.Marshal(struct {
		Tool    string              `json:"tool"`
		Results []map[string]string `json:"results"`
	}{
		Tool:    msfExecutable,
		Results: results,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(respBytes), nil
}

func GenerateMsfQueries(ports []PortInfo) []string {
	var queries []string
	for _, port := range ports {
		if strings.ToLower(port.State) == "open" {
			queries = append(queries, fmt.Sprintf("type:exploit name:%s", port.Service))
			queries = append(queries, fmt.Sprintf("port %s", port.Port))
		}
	}

	return queries
}

func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			if cfg.MsfDisable {
				return nil, nil
			}

			tool := MsfSearchTool{
				executable:   "",
				argsTemplate: "",
			}

			return &tools.ToolData{
				Definition: MsfSearchToolDefinition,
				Call:       tool.Call,
			}, nil
		},
	)
}
