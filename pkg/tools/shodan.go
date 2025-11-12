// ────────────────────────────────────────────────────────────────────────
//
//	shodan.go – wrapper around github.com/shadowscatcher/shodan
//
//	The tool exposes two top‑level JSON keys: `query` and `params`.
//	- `query`   holds all search filters (product, hostname, org, …).
//	- `params`  holds the page number (and any other flags you may add later).
//
//	It can be called from the executor with a payload such as:
//
//	    shodanQuery := tools.ShodanToolArgs{
//	        Query: search.Query{Hostname:"rzd.ru", Port: []int{22,3333}},
//	        Params: search.Params{Page:1},
//	    }
//
// ────────────────────────────────────────────────────────────────────────
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// The official shodan client (the one you pasted in the question).
	shodanclient "github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/search"

	// Internal executor plumbing
	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"

	// The LLM definition type (the one you used for ReWOOTool, etc.).
	"github.com/tmc/langchaingo/llms"
)

// -------------------------------------------------------------------
// 1️⃣  Definition – tells the executor what keys we expect.
// -------------------------------------------------------------------
var ShodanToolDefinition = llms.FunctionDefinition{
	Name:        "shodan",
	Description: `Builds a shodan search query and returns the raw result as a string.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			// The actual search term
			"query": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"domain":   map[string]any{"type": "string"},
					"product":  map[string]any{"type": "string"},
					"asn":      map[string]any{"type": "string"},
					"hostname": map[string]any{"type": "string"},
					"org":      map[string]any{"type": "string"},
					"net":      map[string]any{"type": "string"},
					"country":  map[string]any{"type": "string"},
					// Port list – a slice of ints
					"port": map[string]any{"type": "array",
						"items": map[string]any{"type": "integer"},
					},
				},
			},

			// Top‑level flags (page, filters…)
			"params": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"page": map[string]any{
						"type":        "uint",
						"description": "Page number of the search (default 1).",
					},
				},
			},
		},
	},
}

// -------------------------------------------------------------------
// 2️⃣  Args struct – matches the JSON we described above.
// -------------------------------------------------------------------
type ShodanToolArgs struct {
	Query  search.Query  `json:"query"`
	Params search.Params `json:"params,omitempty"`
}

// -------------------------------------------------------------------
// 3️⃣  Tool type (empty; only used for method receiver).
// -------------------------------------------------------------------
type ShodanTool struct{}

// -------------------------------------------------------------------
// 4️⃣  Call – receives the JSON payload, runs the shodan client and returns a string.
// -------------------------------------------------------------------
func (s ShodanTool) Call(ctx context.Context, input string) (string, error) {
	// 1. Unmarshal incoming JSON into our args struct
	shodanArgs := ShodanToolArgs{}
	if err := json.Unmarshal([]byte(input), &shodanArgs); err != nil {
		return "", fmt.Errorf("json unmarshal: %w", err)
	}

	// 2. Build a client – the library expects an API key in an env var.
	client, _ := shodanclient.GetClient(
		os.Getenv("SHODAN_API_KEY"), // <-- your key
		http.DefaultClient,
		true, // use API v1
	)
	//ctx = context.Background()

	// 3. Build a search request from the supplied args.
	shodanSearch := search.Params{
		Page:  shodanArgs.Params.Page,
		Query: shodanArgs.Query,
	}

	result, err := client.Search(ctx, shodanSearch)
	if err != nil {
		return "", fmt.Errorf("shodan search failed: %w", err)
	}

	// 4. Return the raw JSON string so that other tools can consume it.
	out, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(out), nil
}

// -------------------------------------------------------------------
// 5️⃣  Register this tool in the global registry (used by NewToolsExecutor).
// -------------------------------------------------------------------
func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			// Skip if you disable it via an env var
			if cfg.ShodanDisable {
				return nil, nil
			}

			shodanTool := ShodanTool{}

			definition := ShodanToolDefinition

			return &tools.ToolData{
				Definition: definition,
				Call:       shodanTool.Call,
			}, nil
		},
	)
}
