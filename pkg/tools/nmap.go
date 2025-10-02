package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"
	"github.com/tmc/langchaingo/llms"
)

var NmapToolDefinition = llms.FunctionDefinition{
	Name:        "nmap",
	Description: "Executes nmap with configurable args (or uses sane defaults), parses the output, and generates Metasploit search queries.",
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ip": map[string]any{
				"type":        "string",
				"description": "The valid IP address to scan.",
			},
			"args": map[string]any{
				"type":        "array",
				"description": "Optional array of nmap arguments (e.g. [\"-sV\",\"-p\",\"1-100\"]). If omitted, defaults are used. Defaults are: [\"-v\",\"-T3\",\"-sT\",\"-sV\",\"-Pn\",\"--version-all\",\"--top-ports\", \"100\",]",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
	},
}

type NmapTool struct{}

type NmapToolArgs struct {
	IP   string   `json:"ip"`
	Args []string `json:"args,omitempty"`
}

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

	if nmapToolArgs.IP == "" || net.ParseIP(nmapToolArgs.IP) == nil {
		return "", fmt.Errorf("invalid or missing IP: %q", nmapToolArgs.IP)
	}

	var args []string
	if len(nmapToolArgs.Args) > 0 {
		args = append([]string{}, nmapToolArgs.Args...)
		ipPresent := false
		for _, a := range args {
			if a == nmapToolArgs.IP {
				ipPresent = true
				break
			}
		}
		if !ipPresent {
			args = append(args, nmapToolArgs.IP)
		}
	} else {
		args = []string{
			"-v",
			"-T3",
			"-sT",
			"-sV",
			"-Pn",
			"--version-all",
			"--top-ports", "100",
			nmapToolArgs.IP,
		}
	}

	cmd := exec.Command("nmap", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute nmap: %w - output: %s", err, string(output))
	}

	ports := ParseNmapPorts(string(output))

	//TODO: should be in another pkg
	msfQueries := GenerateMsfQueries(ports)

	response := fmt.Sprintf("Used args: %v\n\nGenerated Metasploit queries: %v", args, msfQueries)
	return response, nil
}

func ParseNmapPorts(nmapOutput string) []PortInfo {
	var ports []PortInfo
	re := regexp.MustCompile(`(\d+)\/tcp\s+(\w+)\s+(.+?)\s*(?:\n|$)`)
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
