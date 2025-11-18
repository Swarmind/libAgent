package tools

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/util"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func TestCommandExecutor(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		CommandExecutorDefinition.Name,
	))
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			t.Fatalf("cleanup: %v", err)
		}
	}()

	result, err := toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "echo \"banana\""
		}`,
	)

	if result != "banana" {
		t.Errorf("result should be \"banana\", got:\n%s", result)
	}

	result, err = toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "echo \"banana\" > test"
		}`,
	)

	if result != "" {
		t.Errorf("result should be \"\", got:\n%s", result)
	}

	result, err = toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "cat test"
		}`,
	)
	if result != "banana" {
		t.Errorf("result should be \"banana\", got:\n%s", result)
	}

	result, err = toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "ls"
		}`,
	)
	if result != "test" {
		t.Errorf("result should be \"test\", got:\n%s", result)
	}
}

func TestCommandExecutorUtils(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		CommandExecutorDefinition.Name,
	))
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			t.Fatalf("cleanup: %v", err)
		}
	}()

	result, err := toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "nmap -F google.com"
		}`,
	)

	if !strings.Contains(result, "443/tcp") {
		t.Errorf("result should contain \"443/tcp\", got:\n%s", result)
	}

	result, err = toolsExecutor.CallTool(
		t.Context(),
		CommandExecutorDefinition.Name,
		`{
			"command": "dig google.com"
		}`,
	)

	if !strings.Contains(result, "google.com.") {
		t.Errorf("result should contain \"google.com.\", got:\n%s", result)
	}
}

func TestCommandExecutorLLMCall(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		CommandExecutorDefinition.Name,
	))
	if err != nil {
		t.Fatalf("new executor: %v", err)
	}
	defer func() {
		if err := toolsExecutor.Cleanup(); err != nil {
			t.Fatalf("cleanup: %v", err)
		}
	}()

	options := []llms.CallOption{
		llms.WithTools([]llms.Tool{
			{
				Type:     "function",
				Function: &CommandExecutorDefinition,
			},
		}),
	}
	llm, err := openai.New(
		openai.WithBaseURL(cfg.AIURL),
		openai.WithToken(cfg.AIToken),
		openai.WithModel(cfg.Model),
		openai.WithAPIVersion("v1"),
	)
	if err != nil {
		t.Fatalf("new llm: %v", err)
	}
	response, err := llm.GenerateContent(t.Context(),
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman,
				"use command executor to echo word 'banana'",
			)},
		options...,
	)
	if err != nil {
		t.Fatalf("generate content: %v", err)
	}
	content := response.Choices[0].Content
	if toolContent := toolsExecutor.ProcessToolCalls(
		t.Context(), response.Choices[0].ToolCalls,
	); toolContent != "" {
		content = toolContent
	}

	jsonSafeContent, err := json.Marshal(util.RemoveThinkTag(content))
	if err != nil {
		t.Fatalf("json safe content: %v", err)
	}

	fmt.Printf("TestCommandExecutorLLMCall result: %s\n", jsonSafeContent)
}
