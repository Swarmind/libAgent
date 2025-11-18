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

func TestWebReaderHTML(t *testing.T) {
	tool := WebReaderTool{}
	result, err := tool.Call(t.Context(), `{
		"url":"https://ifconfig.me"
	}`)
	if err != nil {
		t.Errorf("tool call: %v", err)
	}
	if !strings.Contains(result, "What Is My IP Address") {
		t.Errorf("result not contains \"What Is My IP Address\": %s", result)
	}
}

func TestWebReaderJSON(t *testing.T) {
	tool := WebReaderTool{}
	result, err := tool.Call(t.Context(), `{
		"url":"https://internetdb.shodan.io/1.1.1.1"
	}`)
	if err != nil {
		t.Errorf("tool call: %v", err)
	}
	if !strings.Contains(result, "cloudflare") {
		t.Errorf("result not contains \"cloudflare\": %s", result)
	}
}

func TestWebReaderJS(t *testing.T) {
	tool := WebReaderTool{}
	result, err := tool.Call(t.Context(), `{
		"url":"https://swarmind.ai"
	}`)
	if err != nil {
		t.Errorf("tool call: %v", err)
	}
	if !strings.Contains(result, "https://swarmind.ai/icons/icon.svg") {
		t.Errorf("result not contains \"https://swarmind.ai/icons/icon.svg\": %s", result)
	}
}

func TestWebReaderLLMCall(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		WebReaderDefinition.Name,
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
				Function: &WebReaderDefinition,
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
				"call ifconfig.me resource to figure out host ip address",
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

	fmt.Printf("TestWebReaderLLMCall result: %s\n", jsonSafeContent)
}
