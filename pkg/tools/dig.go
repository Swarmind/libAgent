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

// 2. Definition – tells the executor what arguments the tool expects.
var DigToolDefinition = llms.FunctionDefinition{
	Name:        "dig",
	Description: `Executes dig on a domain name, parses the IP and returns it as JSON.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"domain": map[string]any{
				"type":        "string",
				"description": "The domain name to resolve (e.g. example.com).",
			},
			// In future you can add “timeout”, “dns‑server” etc.
		},
	},
}

// 3. Args struct – matches the JSON keys of the definition above
type DigArgs struct {
	Domain string `json:"domain"`
}

// 4. Tool implementation – we use exec.Command like in NmapTool.
type DigTool struct{}

// Call is invoked by the executor (see Example 1)
func (s DigTool) Call(ctx context.Context, input string) (string, error) {
	// Unmarshal the JSON payload into our args struct
	digArgs := DigArgs{}
	if err := json.Unmarshal([]byte(input), &digArgs); err != nil {
		return "", fmt.Errorf("unmarshal dig: %w", err)
	}

	// Build and run the command.  We use -x to get a single line per answer.
	cmd := exec.Command("dig", "-x", digArgs.Domain, "+short")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("executing dig: %w – output: %s", err, string(out))
	}

	// The raw output from “dig -x … +short” is simply the IP address.
	ip := strings.TrimSpace(string(out))

	// Return a clean string so that other tools can consume it
	return fmt.Sprintf("IP for %s = %s", digArgs.Domain, ip), nil
}

// 5. Register with the globalToolsRegistry (see Example 4 init)
func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			// Skip if you added a “digDisable” flag in your config
			if cfg.DigDisable { // <-- optional – see below
				return nil, nil
			}

			return &tools.ToolData{
				Definition: DigToolDefinition,
				Call:       DigTool{}.Call,
			}, nil
		},
	)
}
