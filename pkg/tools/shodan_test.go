package tools

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/util"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func TestShodanSimplePayload(t *testing.T) {
	tool := ShodanTool{}
	params, err := tool.UnmarshalFromString(
		`{
		"params": {
			"query": {
				"ip": "127.0.0.1",
				"has_ssl": true,
				"bitcoin": {
					"ip": "127.0.0.1"
				}
			},
			"page": 2
		}
	}`,
	)
	if err != nil {
		t.Errorf("unmarshal: %v", err)
	}
	if params.Query.IP != "127.0.0.1" {
		t.Errorf("params.Query.IP != 127.0.0.1, got: %s", params.Query.IP)
	}
	if params.Query.HasSSL != true {
		t.Errorf("params.Query.HasSSL != true, got: %t", params.Query.HasSSL)
	}
	if params.Query.Bitcoin.IP != "127.0.0.1" {
		t.Errorf("params.Query.Bitcoin.IP != 127.0.0.1, got: %s", params.Query.Bitcoin.IP)
	}
	if params.Page != 2 {
		t.Errorf("params.Page != 2, got: %d", params.Page)
	}
}

func TestShodanFullPayload(t *testing.T) {
	tool := ShodanTool{}
	params, err := tool.UnmarshalFromString(
		`{"params":{"facets":[],"minify":true,"offset":0,"page":1,"query":{"Text":"127.0.0.1","after":"","all":"","asn":"","before":"","bitcoin":{"ip":"","ip_count":0,"port":0,"version":""},"city":"","country":"","cpe":"","device":"","geo":"","has_ipv6":false,"has_screenshot":false,"has_ssl":false,"has_vuln":false,"hash":0,"hostname":"","http":{"component":"","component_category":"","favicon":{"hash":0},"html":"","html_hash":0,"robots_hash":0,"securitytxt":"","status":0,"title":"","waf":""},"ip":"127.0.0.1","isp":"","link":"","net":"","ntp":{"ip":"","ip_count":0,"more":false,"port":0},"org":"","os":"","port":0,"postal":"","product":"","region":"","scan":"","screenshot":{"label":""},"shodan":{"module":""},"snmp":{"contact":"","location":"","name":""},"ssh":{"hassh":"","type":""},"ssl":{"alpn":"","cert":{"alg":"","expired":false,"extension":"","fingerprint":"","issuer":{"cn":""},"pubkey":{"bits":0,"type":""},"serial":"","subject":{"cn":""}},"chain_count":0,"cipher":{"bits":0,"name":"","version":""},"version":""},"state":"","tag":"","telnet":{"do":"","dont":"","option":"","will":"","wont":""},"version":"","vuln":""}}}`,
	)
	if err != nil {
		t.Errorf("unmarshal: %v", err)
	}
	if params.Query.Text != "127.0.0.1" {
		t.Errorf("params.Query.Text != 127.0.0.1, got: %s", params.Query.Text)
	}
	if params.Query.IP != "127.0.0.1" {
		t.Errorf("params.Query.IP != 127.0.0.1, got: %s", params.Query.IP)
	}
}

func TestShodanLLMCall(t *testing.T) {
	cfg, err := config.NewConfig("../../.env")
	if err != nil {
		t.Fatalf("new config: %v", err)
	}

	toolsExecutor, err := NewToolsExecutor(t.Context(), cfg, WithToolsWhitelist(
		ShodanToolDefinition.Name,
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
				Function: &ShodanToolDefinition,
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
				"call shodan tool over 127.0.0.1 ip",
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

	fmt.Printf("TestShodanLLMCall result: %s\n", jsonSafeContent)
}
