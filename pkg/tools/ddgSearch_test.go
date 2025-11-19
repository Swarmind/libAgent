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
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

func TestDDGSearch(t *testing.T) {
	wrappedTool, err := duckduckgo.New(
		5,
		duckduckgo.DefaultUserAgent,
	)
	if err != nil {
		t.Fatalf("duckduckgo new: %v", err)
	}
	tool := DDGSearchTool{
		wrappedTool: wrappedTool,
	}

	result, err := tool.Call(t.Context(), `{
		"query":"bananas"
	}`)

	if c := strings.Count(result, "Title: "); c != 5 {
		t.Errorf("result should contain 5 \"Title: \", not %d: %s", c, result)
	}
	if c := strings.Count(result, "Description: "); c != 5 {
		t.Errorf("result should contain 5 \"Description: \", not %d: %s", c, result)
	}
	if c := strings.Count(result, "URL: "); c != 5 {
		t.Errorf("result should contain 5 \"URL: \", not %d: %s", c, result)
	}
}

func TestDDGSearchLLMCall(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		DDGSearchDefinition.Name,
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
				Function: &DDGSearchDefinition,
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
				"call ddg search tool and search for 'banana'",
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

	fmt.Printf("TestDDGSearchLLMCall result: %s\n", jsonSafeContent)
}
