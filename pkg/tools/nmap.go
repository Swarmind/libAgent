package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/tmc/langchaingo/llms"
)

var NmapToolDefinition = llms.FunctionDefinition{
	Name:        "nmap",
	Description: "Executes nmap -v -T4 -PA -sV --version-all --osscan-guess -A -sS -p 1-65535 [IP], parses the output, and generates Metasploit search queries.",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ip": map[string]any{
				"type":        "string",
				"description": "The valid IP address to scan to.",
			},
		},
	},
}

type NmapTool struct{}

type NmapToolArgs struct {
	IP string `json:"ip"`
}

// struct для описания найденного порта
type PortInfo struct {
	Port    string
	State   string
	Service string
}

// Call executes the command with the given arguments.
func (s NmapTool) Call(ctx context.Context, input string) (string, error) {
	nmapToolArgs := NmapToolArgs{}
	if err := json.Unmarshal([]byte(input), &nmapToolArgs); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	cmd := exec.Command("nmap", "-v", "-T4", "-PA", "-sV", "--version-all", "-osscan-guess", "-A", "-sS", "-p", "1-65535", nmapToolArgs.IP)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute Nmap: %w", err)
	}

	// Парсинг результатов Nmap
	ports := parseNmapPorts(string(output))
	msfQueries := generateMsfQueries(ports)

	// Формирование ответа для LangChainGo
	response := fmt.Sprintf("Generated Metasploit queries: %v", msfQueries)
	return response, nil
}

func parseNmapPorts(nmapOutput string) []PortInfo {
	var ports []PortInfo
	re := regexp.MustCompile(`(\d+)\/tcp\s+(\w+)\s+(.+?)\s*`)
	matches := re.FindAllStringSubmatch(nmapOutput, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			ports = append(ports, PortInfo{
				Port:    match[1],
				State:   strings.TrimSpace(match[2]),
				Service: strings.TrimSpace(match[3]),
			})
		}
	}

	return ports
}

func generateMsfQueries(ports []PortInfo) []string {
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
			if cfg.NmapDisable {
				return nil, nil
			}

			return &tools.ToolData{
				Definition: NmapToolDefinition,
				Call:       NmapTool{}.Call,
			}, nil
		},
	)
}
