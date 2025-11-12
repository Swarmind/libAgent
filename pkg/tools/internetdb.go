// ────────────────────────────────────────────────────────────────────────
//
//	internetdb.go – wrapper around https://internetdb.shodan.io
//
//	The tool simply performs a GET request to https://internetdb.shodan.io/<ip>
//	(e.g. https://internetdb.shodan.io/1.1.1.1) and returns the raw JSON.
//
//	Usage:
//
//	    dbQuery := tools.InternetDBArgs{IP:"8.8.4.4"}
//	    dbResult, err := toolsExecutor.CallTool(ctx,
//	        tools.InternetDBToolDefinition.Name,
//	        string(json.Marshal(dbQuery)),
//	    )
//
// ────────────────────────────────────────────────────────────────────────
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// Internal executor plumbing
	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"

	// The LLM definition type (the one you used for ReWOOTool, etc.).
	"github.com/tmc/langchaingo/llms"
)

// -------------------------------------------------------------------
// 1️⃣  Definition – tells the executor what keys we expect.
// -------------------------------------------------------------------
var InternetDBToolDefinition = llms.FunctionDefinition{
	Name:        "internetdb",
	Description: `Executes a GET request to https://internetdb.shodan.io/<ip> and returns raw JSON.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			// Only one property – the IP address.
			"query": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"ip": map[string]any{"type": "string"},
				},
			},
		},
	},
}

// -------------------------------------------------------------------
// 2️⃣  Args struct – matches the JSON we described above.
// -------------------------------------------------------------------
type InternetDBArgs struct {
	Query struct {
		IP string `json:"ip"`
	} `json:"query"`
}

// -------------------------------------------------------------------
// 3️⃣  Tool type (empty; only used for method receiver).
// -------------------------------------------------------------------
type InternetDBTool struct{}

// -------------------------------------------------------------------
// 4️⃣  Call – receives the JSON payload, performs GET and returns a string.
// -------------------------------------------------------------------
func (s *InternetDBTool) Call(ctx context.Context, input string) (string, error) {
	// 1. Unmarshal incoming JSON into our args struct.
	var inArgs InternetDBArgs
	if err := json.Unmarshal([]byte(input), &inArgs); err != nil {
		return "", fmt.Errorf("json unmarshal: %w", err)
	}

	// 2. Build the request URL.
	url := fmt.Sprintf("https://internetdb.shodan.io/%s", inArgs.Query.IP)

	// 3. Execute a GET and read whole body.
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP get: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	// 4. Return the raw JSON string.
	return string(body), nil
}

// -------------------------------------------------------------------
// 5️⃣  Register this tool in the global registry (used by NewToolsExecutor).
// -------------------------------------------------------------------
func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			// Skip if you disable it via an env var
			if cfg.InternetDBDisable {
				return nil, nil
			}

			dbTool := InternetDBTool{}
			definition := InternetDBToolDefinition

			return &tools.ToolData{
				Definition: definition,
				Call:       dbTool.Call,
			}, nil
		},
	)
}
