package tools

import (
	"context"
	"encoding/json"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/internal/tools/rewoo"
	"github.com/Swarmind/libagent/pkg/config"

	graph "github.com/JackBekket/langgraphgo/graph/stategraph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

var ReWOOToolDefinition = llms.FunctionDefinition{
	Name: "rewoo",
	Description: `A more complex LLM Reasoning algorithm.
Useful when you need to do a tool assisted reasoning research.
Usually tends to return a short response as a result of multiple step thinking.
Use it is you think that you have an isolated complex research subtask.
Input can be any complex task.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "The task query",
			},
		},
	},
}

type ReWOOToolArgs struct {
	Query string `json:"query"`
}

type ReWOOTool struct {
	ReWOO rewoo.ReWOO
	graph *graph.Runnable
}

func (t *ReWOOTool) Call(ctx context.Context, input string) (string, error) {
	rewooToolArgs := ReWOOToolArgs{}
	if err := json.Unmarshal([]byte(input), &rewooToolArgs); err != nil {
		return "", err
	}

	if t.ReWOO.ToolsExecutor == nil {
		t.ReWOO.ToolsExecutor = globalToolsExecutor
	}
	if t.graph == nil {
		if g, err := t.ReWOO.InitializeGraph(); err != nil {
			return "", err
		} else {
			t.graph = g
		}
	}

	state, err := t.graph.Invoke(ctx, &rewoo.State{
		Task: rewooToolArgs.Query,
	})
	if err != nil {
		return "", err
	}

	return state.(*rewoo.State).Result, nil
}

func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			if cfg.ReWOODisable {
				return nil, nil
			}
			llm, err := openai.New(
				openai.WithBaseURL(cfg.AIURL),
				openai.WithToken(cfg.AIToken),
				openai.WithModel(cfg.Model),
				openai.WithAPIVersion("v1"),
			)
			if err != nil {
				return nil, err
			}

			rewooTool := ReWOOTool{
				ReWOO: rewoo.ReWOO{
					LLM:                llm,
					DefaultCallOptions: config.ConifgToCallOptions(cfg.ReWOODefaultCallOptions),
					MaxAttempts:        cfg.ReWOOMaxAttempts,
				},
			}

			return &tools.ToolData{
				Definition: ReWOOToolDefinition,
				Call:       rewooTool.Call,
			}, nil
		},
	)
}
