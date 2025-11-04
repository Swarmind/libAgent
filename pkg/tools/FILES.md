# pkg/tools/attacker.go  
## Package / Component    
**tools**  
  
### Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides the `Context` type used in the tool’s `Call` method. |  
| `encoding/json` | Handles JSON unmarshalling of input arguments. |  
| `fmt` | Used for printing debug information to stdout. |  
| `os/exec` | Executes external commands (`msfconsole`). |  
| `github.com/Swarmind/libagent/internal/tools` | Provides the `ToolData` struct and a global registry used in `init`. |  
| `github.com/Swarmind/libagent/pkg/config` | Supplies configuration data for enabling/disabling the tool. |  
| `github.com/tmc/langchaingo/llms` | Supplies the LLM function definition type (`FunctionDefinition`). |  
  
### External Data / Input Sources    
* The tool expects a JSON string containing two fields:    
  * `"module"` – a string specifying the Metasploit module to use (e.g., `"exploit/module_name"`).    
  * `"options"` – an object of key‑value pairs for Metasploit options (e.g., `{"RHOSTS":"192.168.1.10","LHOST":"192.168.1.5","payload":"cmd/unix/reverse"}`).    
  
The tool will parse this JSON into a Go struct, build command arguments for `msfconsole`, execute the exploit, and return the raw output.  
  
### TODOs    
No explicit TODO comments are present in the file; all functionality is already implemented.  
  
---  
  
## Major Code Parts  
  
### 1. `ExploitToolDefinition` (LLM function definition)    
Defines a function named **"exploit"** that can be invoked by an LLM. The JSON schema describes two properties:    
* `module`: string – the Metasploit module to run.    
* `options`: object – arbitrary key/value options for the exploit.  
  
### 2. `ExploitToolArgs` struct    
Holds the parsed input data:  
```go  
type ExploitToolArgs struct {  
    Module  string            `json:"module"`  
    Options map[string]string `json:"options"`  
}  
```  
The struct tags match the JSON schema defined above.  
  
### 3. `(*ExploitTool).Call` method    
* Unmarshals the input JSON into an `ExploitToolArgs`.    
* Builds a command slice for running `msfconsole use cmd/unix/reverse`.    
* Executes that command and captures its output.    
* Builds another command slice to run the user‑specified module with all provided options (`-o key=value`).    
* Executes the second command, prints debug information, and returns the raw output string.  
  
### 4. `init` function    
Registers the tool in a global registry so that it can be discovered by other components. It checks the configuration flag `cfg.ExploitDisable`; if false, it creates an instance of `ExploitTool`, wraps it into a `tools.ToolData` struct (containing the definition and call method), and appends it to `globalToolsRegistry`.  
  
---  
  
This file implements a single LLM‑exposed tool that runs a Metasploit exploit via `msfconsole`. It is ready for integration into the larger package, providing both configuration handling and runtime execution logic.  
  
# pkg/tools/ddgSearch.go  
# Package / Component    
**tools**  
  
## Imports  
```go  
import (  
	"context"  
	"encoding/json"  
  
	"github.com/Swarmind/libagent/internal/tools"  
	"github.com/Swarmind/libagent/pkg/config"  
  
	"github.com/tmc/langchaingo/llms"  
	"github.com/tmc/langchaingo/tools/duckduckgo"  
)  
```  
* `context` – standard Go context handling.    
* `encoding/json` – JSON marshaling/unmarshaling for input arguments.    
* `github.com/Swarmind/libagent/internal/tools` – internal tool registry and data structures.    
* `github.com/Swarmind/libagent/pkg/config` – configuration struct used by the init routine.    
* `github.com/tmc/langchaingo/llms` – LLM function definition type.    
* `github.com/tmc/langchaingo/tools/duckduckgo` – DuckDuckGo search tool implementation.  
  
## External Data / Input Sources  
| Config field | Purpose | Default value |  
|--------------|---------|---------------|  
| `cfg.DDGSearchDisable` | Flag to skip registration of the tool | N/A (used as a guard) |  
| `cfg.DDGSearchMaxResults` | Max number of results returned by DuckDuckGo | 5 (default if zero) |  
| `cfg.DDGSearchUserAgent` | User‑agent string for DuckDuckGo requests | `duckduckgo.DefaultUserAgent` |  
  
## TODOs  
No explicit TODO comments are present in the file.  
  
## Summary of Major Code Parts  
  
### DDGSearchDefinition    
A global variable of type `llms.FunctionDefinition`. It declares a function named **webSearch** with a single string parameter `query`. The definition is used by the LLM to understand what arguments the tool expects and how it should be called.  
  
```go  
var DDGSearchDefinition = llms.FunctionDefinition{  
	Name: "webSearch",  
	Description: `A duckduckgo search wrapper.  
Given search query returns a multiple results with short descriptions and URLs.`,  
	Parameters: map[string]any{ ... },  
}  
```  
  
### DDGSearchArgs    
Simple struct that matches the JSON payload expected by the tool. It contains one field, `Query`, which is tagged for JSON unmarshaling.  
  
```go  
type DDGSearchArgs struct {  
	Query string `json:"query"`  
}  
```  
  
### DDGSearchTool & Call Method    
`DDGSearchTool` holds a reference to the underlying DuckDuckGo search implementation (`*duckduckgo.Tool`). The `Call` method receives a context and an input string, unmarshals it into `DDGSearchArgs`, then forwards the query to the wrapped tool.  
  
```go  
type DDGSearchTool struct {  
	wrappedTool *duckduckgo.Tool  
}  
  
func (t DDGSearchTool) Call(ctx context.Context, input string) (string, error) {  
	ddgSearchArgs := DDGSearchArgs{}  
	if err := json.Unmarshal([]byte(input), &ddgSearchArgs); err != nil {  
		return "", err  
	}  
	return t.wrappedTool.Call(ctx, ddgSearchArgs.Query)  
}  
```  
  
### init Function    
Registers the tool in a global registry (`globalToolsRegistry`). It checks configuration flags, applies defaults for max results and user agent, creates a DuckDuckGo instance via `duckduckgo.New`, wraps it into a `DDGSearchTool`, and finally returns a `*tools.ToolData` containing the function definition and the Call method.  
  
```go  
func init() {  
	globalToolsRegistry = append(globalToolsRegistry,  
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {  
			if cfg.DDGSearchDisable {  
				return nil, nil  
			}  
			if cfg.DDGSearchMaxResults == 0 {  
				cfg.DDGSearchMaxResults = 5  
			}  
			if cfg.DDGSearchUserAgent == "" {  
				cfg.DDGSearchUserAgent = duckduckgo.DefaultUserAgent  
			}  
  
			wrappedTool, err := duckduckgo.New(  
				cfg.DDGSearchMaxResults,  
				cfg.DDGSearchUserAgent,  
			)  
			if err != nil {  
				return nil, err  
			}  
  
			ddgSearchTool := DDGSearchTool{  
				wrappedTool: wrappedTool,  
			}  
  
			return &tools.ToolData{  
				Definition: DDGSearchDefinition,  
				Call:       ddgSearchTool.Call,  
			}, nil  
		},  
	)  
}  
```  
  
The file thus provides a fully‑functional DuckDuckGo search wrapper that can be invoked by the LLM system, with configuration-driven defaults and registration into the global tool registry.  
  
# pkg/tools/executor.go  
# Package / Component    
**tools**  
  
## Imports  
```go  
import (  
	"context"  
	"encoding/json"  
	"fmt"  
	"os"  
	"strings"  
	"time"  
  
	"github.com/Swarmind/libagent/internal/tools"  
	"github.com/Swarmind/libagent/pkg/config"  
	"github.com/ThomasRooney/gexpect"  
	"github.com/google/uuid"  
	"github.com/rs/zerolog/log"  
  
	"github.com/tmc/langchaingo/llms"  
)  
```  
  
## External Data / Input Sources  
| Source | Description |  
|--------|-------------|  
| `config.Config` | Provides configuration flags and command lists (`CommandExecutorDisable`, `CommandExecutorCommands`). |  
| `github.com/tmc/langchaingo/llms.FunctionDefinition` | Used to describe the tool’s function signature. |  
| `gexpect.ExpectSubprocess` | Handles interactive bash sessions. |  
  
## TODOs  
No explicit TODO comments were found in this file.  
  
---  
  
# Summary of Major Code Parts  
  
## 1. Constants & Variable Definitions  
- **CommandNotesPromptAddition** – a string that will be appended to the tool description when registering.  
- **sendChunkSize** – buffer size for sending commands (currently unused but defined).  
- **CommandExecutorDefinition** – an `llms.FunctionDefinition` instance that describes the tool’s name, description, and JSON schema.    
  *Name*: `"commandExecutor"`    
  *Description*: “Executes a provided command in a interactive stateful bash shell session…”.    
  *Parameters*: expects a single string field called `command`.  
  
## 2. Data Structures  
### `CommandExecutorArgs`  
```go  
type CommandExecutorArgs struct {  
	Command string `json:"command"`  
}  
```  
Used to unmarshal the JSON payload passed to the tool.  
  
### `CommandExecutorTool`  
```go  
type CommandExecutorTool struct {  
	tempDir *string  
	process *gexpect.ExpectSubprocess  
	prompt  string  
}  
```  
Holds state for a single execution session: temporary directory, subprocess handle, and a unique prompt string.  
  
## 3. Methods  
  
### `Call(ctx context.Context, input string) (string, error)`  
- Unmarshals the JSON payload into `CommandExecutorArgs`.  
- Delegates to `RunCommand` with the extracted command string.  
- Returns the raw output of the executed command.  
  
### `RunCommand(input string) (string, error)`  
1. **Initialization**    
   - If no temp directory exists, create one (`os.MkdirTemp`) and spawn a bash process in it via `gexpect.SpawnAtDirectory`.    
   - Configure environment variables (`PATH`, `HOME`, `GOCACHE`), set a unique prompt using a UUID, and prepare the shell for raw mode.    
2. **Command Execution**    
   - Trim trailing backslash from input, append newline, send it to the subprocess, wait for the custom prompt, collect output.  
3. **Return**    
   - Returns trimmed command output as string.  
  
### `cleanup() error`  
- Removes the temporary directory and closes the subprocess if a temp directory was created.  
  
## 4. Tool Registration (`init` function)  
- Appends a closure to `globalToolsRegistry`.    
- The closure:  
  * Checks `cfg.CommandExecutorDisable`; returns nil if disabled.  
  * Builds a list of commands from `cfg.CommandExecutorCommands`, formatting each as “- command: description”.  
  * Extends the tool definition’s description with this list and the constant prompt addition.  
  * Returns a populated `tools.ToolData` containing the definition, Call, and Cleanup functions.  
  
---  
  
This file defines a reusable “command executor” tool that can be invoked by other components in the package. It handles environment setup, command execution, and cleanup, while exposing its configuration through the global registry.  
  
# pkg/tools/metasploit.go  
## Package / Component    
**tools**  
  
### Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides the context type used in callbacks and registrations. |  
| `encoding/json` | Marshal/unmarshal JSON payloads for tool input & output. |  
| `fmt` | String formatting, printing, and error handling. |  
| `os/exec` | Execute external commands (Metasploit console). |  
| `strings` | Lower‑casing port state strings in query generation. |  
| `github.com/Swarmind/libagent/internal/tools` | Tool registry & data structures (`ToolData`). |  
| `github.com/Swarmind/libagent/pkg/config` | Configuration struct used during registration. |  
| `github.com/tmc/langchaingo/llms` | LLM function definition type (`FunctionDefinition`). |  
  
---  
  
## External Data / Input Sources  
* **Config** – `config.Config` is passed to the init callback.  
* **PortInfo** – slice of ports supplied to `GenerateMsfQueries`. (Type defined elsewhere in the project.)  
* **JSON input string** – expected to contain a field `"queries"` with an array of strings.  
  
---  
  
## TODOs  
No explicit `TODO:` comments are present in this file.    
(If future work is needed, add them here.)  
  
---  
  
## Summary of Major Code Parts  
  
### 1. Function Definition (`MsfSearchToolDefinition`)  
Defines the LLM function “msf_search” with a JSON schema for an array of search queries. The description lists all possible Metasploit query keywords.  
  
### 2. Tool Structs  
* **`MsfSearchTool`** – holds executable name and argument template.  
* **`MsfSearchToolArgs`** – struct used to unmarshal the incoming JSON payload (`queries []string`).  
  
### 3. Global Variables & Setter  
```go  
var (  
    msfExecutable   = "msfconsole"  
    msfArgsTemplate = "search %s; exit"  
)  
func SetMsfCommand(executable string, argsTemplate string) { … }  
```  
Allows overriding the default executable and argument template.  
  
### 4. Call Method (`(s MsfSearchTool).Call`)  
* Unmarshals input JSON into `MsfSearchToolArgs`.  
* Iterates over each query:  
  * Builds command arguments: `-q -x <formatted query>`.  
  * Executes via `exec.Command` and collects output.  
* Aggregates results into a slice of maps (`query`, `output`) and marshals them back to JSON for the LLM.  
  
### 5. Query Generator (`GenerateMsfQueries`)  
Takes a slice of `PortInfo` and produces two queries per open port:  
1. `"type:exploit name:<service>"`  
2. `"port <port>"`  
  
### 6. Registration (`init()`)  
Appends a registration function to the global tool registry, creating an instance of `MsfSearchTool` (with defaults) and returning its definition & call method.  
  
---  
  
All parts together provide a reusable Metasploit search tool that can be invoked by an LLM pipeline, with dynamic query generation from port data.  
  
# pkg/tools/nmap.go  
**Package / Component**    
`tools`  
  
---  
  
### Imports  
```go  
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
```  
  
### External data / input sources  
| Source | Type | Notes |  
|--------|------|-------|  
| `context.Context` | passed to `Call` | Execution context for the tool |  
| JSON string (`input`) | used in `Call` | Contains IP and optional args |  
| `config.Config` | accessed in `init()` | Configuration flag `NmapDisable` |  
| `net.ParseIP` | validates IP address | |  
| `exec.Command` | runs external `nmap` binary | |  
| `regexp.MustCompile` | regex for parsing nmap output | |  
  
### TODOs  
- **Call**: `//TODO: should be in another pkg` – the generation of Metasploit queries (`GenerateMsfQueries`) is currently inline but could be extracted.  
  
---  
  
## Summary of major code parts  
  
#### 1. `NmapToolDefinition`  
Defines a LangChainGo function definition for an “nmap” tool, including:  
- Name and description.  
- Parameters: IP (string) and optional args array (array of strings).  
- Used by the internal registry to expose the tool.  
  
#### 2. Data structures  
```go  
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
```  
- `NmapTool` is a marker type for the tool implementation.  
- `NmapToolArgs` maps JSON input to Go fields.  
- `PortInfo` holds parsed port data.  
  
#### 3. `Call` method  
Executes nmap with supplied or default arguments, parses output, and returns a formatted string:  
1. Unmarshals JSON into `nmapToolArgs`.  
2. Validates the IP address.  
3. Builds an argument slice: if user provided args, ensures IP is included; otherwise uses sane defaults (`-v`, `-F`, etc.).  
4. Runs `exec.Command("nmap", args...)` and captures combined output.  
5. Calls `ParseNmapPorts` to extract port information.  
6. Generates Metasploit queries via `GenerateMsfQueries(ports)` (TODO: move elsewhere).  
7. Returns a string summarizing used arguments and generated queries.  
  
#### 4. `ParseNmapPorts`  
Parses the raw nmap output using regex:  
- Compiles pattern `(\d+)\/tcp\s+(\w+)\s+(.+?)\s*(?:\n|$)`.  
- Finds all matches, appends each as a `PortInfo` entry.  
- Prints debug logs for each found port and final list.  
  
#### 5. `init()` registration  
Adds the tool to a global registry (`globalToolsRegistry`) so that it can be invoked by other components:  
```go  
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
```  
The closure checks a config flag (`NmapDisable`) before registering.  
  
---  
  
# pkg/tools/rewoo.go  
## Package / Component    
**Name:** `tools`    
  
### Imports    
| Import Path | Alias (if any) |  
|-------------|-----------------|  
| `context` | – |  
| `encoding/json` | – |  
| `github.com/Swarmind/libagent/internal/tools` | – |  
| `github.com/Swarmind/libagent/internal/tools/rewoo` | – |  
| `github.com/Swarmind/libagent/pkg/config` | – |  
| `graph "github.com/JackBekket/langgraphgo/graph/stategraph"` | alias `graph` |  
| `github.com/tmc/langchaingo/llms` | – |  
| `github.com/tmc/langchaingo/llms/openai` | – |  
  
### External Data / Input Sources    
* **Config** – The tool is configured via a `config.Config` struct (fields used: `ReWOODisable`, `AIURL`, `AIToken`, `Model`, `RewOODefaultCallOptions`).    
* **LLM** – An OpenAI client created with the above config values.    
* **rewoo.ReWOO** – Holds the LLM and default call options; used to initialize a state‑graph.    
* **graph.Runnable** – The runnable graph that executes the reasoning steps.    
  
### TODOs    
No explicit `TODO:` comments are present in this file, but future improvements could include:    
1. Add error handling for JSON unmarshalling failures.    
2. Cache or reuse the LLM client instead of creating it on every init call.    
  
---  
  
## Summary of Major Code Parts  
  
### 1. Tool Definition (`ReWOOToolDefinition`)    
A global variable of type `llms.FunctionDefinition` that describes the tool to the LLM system:    
* **Name** – `"rewoo"`    
* **Description** – A short explanation of its purpose and expected input format.    
* **Parameters** – JSON schema for a single string field `"query"`.    
  
This definition is later used when registering the tool in `init()`.  
  
### 2. Argument Struct (`ReWOOToolArgs`)    
Defines the shape of the JSON payload that will be passed to the tool’s `Call` method:    
```go  
type ReWOOToolArgs struct {  
    Query string `json:"query"`  
}  
```  
The tag ensures proper unmarshalling from the incoming JSON.  
  
### 3. Tool Implementation (`ReWOOTool`)    
* **Fields** –    
  * `ReWOO rewoo.ReWOO` – holds the LLM client and default options.    
  * `graph *graph.Runnable` – a pointer to the state‑graph that will be executed.    
  
#### Call Method    
```go  
func (t *ReWOOTool) Call(ctx context.Context, input string) (string, error)  
```  
1. Unmarshals the incoming JSON into `rewooToolArgs`.    
2. Lazily creates an OpenAI client if none is set (`t.ReWOO.ToolsExecutor`).    
3. Initializes the graph once (`t.graph`) via `t.ReWOO.InitializeGraph()`.    
4. Invokes the graph with a new state containing the query string, and returns the resulting result string.  
  
### 4. Registration Function (`init()`)    
Registers this tool in a global registry so that it can be discovered by other components:    
  
* Checks if the config flag `ReWOODisable` is true.    
* Creates an OpenAI client with base URL, token, model, and API version from the config.    
* Builds a `ReWOOTool` instance with the LLM and default call options derived from the config.    
* Appends a closure to `globalToolsRegistry` that returns a `tools.ToolData` containing the definition and the Call method.  
  
---  
  
This file defines a reusable “rewoo” tool for complex reasoning tasks, wiring together an OpenAI client, a state‑graph executor, and configuration data into a single component.  
  
# pkg/tools/semanticSearch.go  
## Package & Imports    
**Package name:** `tools`    
  
Imports used in this file:    
  
| Import path | Purpose |  
|-------------|---------|  
| `context` | Provides context handling for DB and LLM calls |  
| `encoding/json` | JSON marshaling/unmarshaling of tool arguments |  
| `fmt` | String formatting for result output |  
| `github.com/Swarmind/libagent/internal/tools` | Tool registry & data structures |  
| `github.com/Swarmind/libagent/pkg/config` | Configuration values for the tool |  
| `github.com/jackc/pgx/v5/pgxpool` | PostgreSQL connection pool |  
| `github.com/tmc/langchaingo/embeddings` | Embedding generation |  
| `github.com/tmc/langchaingo/llms` | LLM interface |  
| `github.com/tmc/langchaingo/llms/openai` | OpenAI client creation |  
| `github.com/tmc/langchaingo/vectorstores/pgvector` | Vector store abstraction |  
  
---  
  
## External Data / Input Sources    
The tool pulls configuration values from a global `config.Config` struct:  
  
- `SemanticSearchDisable` – flag to enable/disable the tool  
- `SemanticSearchAIURL` – base URL for OpenAI LLM  
- `SemanticSearchAIToken` – token for authentication  
- `SemanticSearchDBConnection` – PostgreSQL connection string  
- `SemanticSearchEmbeddingModel` – name of embedding model to use  
- `SemanticSearchMaxResults` – number of results to return  
  
These values are read in the `init()` function and used to instantiate a `SemanticSearchTool`.  
  
---  
  
## TODOs    
| Line | Comment |  
|------|---------|  
| 12 | “there should NOT exist arguments which called NAME cause it cause COLLISION with actual function name. .....more like confusion then collision so there are no error” – clarifies the `collection` parameter naming in the JSON schema |  
  
---  
  
## SemanticSearchDefinition    
A global variable of type `llms.FunctionDefinition`. It declares:  
  
- **Name:** `"semanticSearch"`  
- **Description:** “Performs semantic search in the vector store of the saved code blobs. Returns matching file contents”  
- **Parameters:** an object with two string properties: `query` (search query) and `collection` (name of collection to search).  
  
---  
  
## SemanticSearchArgs    
A simple struct used for unmarshaling JSON input:  
  
```go  
type SemanticSearchArgs struct {  
    Query      string `json:"query"`  
    Collection string `json:"collection"`  
}  
```  
  
It holds the user‑supplied query and target collection name.  
  
---  
  
## SemanticSearchTool    
The core tool configuration:  
  
```go  
type SemanticSearchTool struct {  
    OpenAIURL      string  
    OpenAIToken    string  
    DBConnection   string  
    EmbeddingModel string  
    MaxResults     int  
}  
```  
  
All fields are populated from the global config in `init()`.  
  
---  
  
## Call Method    
`func (s SemanticSearchTool) Call(ctx context.Context, input string) (string, error)`    
  
1. Unmarshals the JSON `input` into a `SemanticSearchArgs`.    
2. Parses the PostgreSQL connection string and creates a pool (`pgxpool`).    
3. Builds an OpenAI LLM client with the configured URL, token, embedding model, and API version `"v1"`.    
4. Creates an embedder from that LLM.    
5. Instantiates a `pgvector` store using the collection name from the arguments, the connection pool, and the embedder.    
6. Executes a similarity search for the query string, limited to `s.MaxResults`.    
7. Concatenates each result’s page content into a single output string.  
  
The method returns that concatenated string along with any error encountered.  
  
---  
  
## init() Registration    
During package initialization:  
  
1. Validates all required config fields and sets defaults where needed (e.g., default max results = 2).    
2. Builds a `SemanticSearchTool` instance from the config values.    
3. Appends a closure to `globalToolsRegistry`. The closure returns a `*tools.ToolData` containing:  
   - `Definition: SemanticSearchDefinition`  
   - `Call: semanticSearchTool.Call`  
  
This makes the tool discoverable by the rest of the application.  
  
---  
  
# pkg/tools/tools.go  
## Package / Component    
**Name:** `tools`    
  
### Imports  
```go  
import (  
	"context"  
	"slices"  
  
	"github.com/Swarmind/libagent/internal/tools"  
	"github.com/Swarmind/libagent/pkg/config"  
)  
```  
* `context` – standard Go context package for request handling.    
* `slices` – standard library helper for slice operations (used to check whitelist membership).    
* `github.com/Swarmind/libagent/internal/tools` – internal tool definitions and executor type.    
* `github.com/Swarmind/libagent/pkg/config` – configuration data used when initializing tools.  
  
---  
  
## External Data / Input Sources  
| Source | Purpose |  
|--------|---------|  
| `context.Context` | Carries request‑level context for tool initialization. |  
| `config.Config` | Configuration values passed to each tool init function. |  
| `globalToolsRegistry` | Slice of functions that create individual tools (`func(context.Context, config.Config) (*tools.ToolData, error)`). |  
  
---  
  
## TODOs  
No explicit TODO comments are present in this file.  
  
---  
  
## Summary of Major Code Parts  
  
### 1. Option & Options Types    
```go  
type ExecutorOption func(*ExecutorOptions)  
  
type ExecutorOptions struct {  
	ToolsWhitelist []string  
}  
```  
* `ExecutorOption` is a functional option that mutates an `ExecutorOptions` instance.    
* `ExecutorOptions` holds configuration for the executor, currently only a whitelist of tool names.  
  
### 2. Global Variables    
```go  
var globalToolsRegistry = []func(context.Context, config.Config) (*tools.ToolData, error){}  
var globalToolsExecutor *tools.ToolsExecutor  
```  
* `globalToolsRegistry` stores init functions that will be called during executor creation.    
* `globalToolsExecutor` keeps a reference to the last created executor for later use.  
  
### 3. `NewToolsExecutor` – Core Factory Function    
```go  
func NewToolsExecutor(ctx context.Context, cfg config.Config, opts ...ExecutorOption) (*tools.ToolsExecutor, error)  
```  
1. **Create local structures** – an empty `toolsExecutor`, a map to hold tool data, and an options instance.    
2. **Apply functional options** – iterate over the variadic `opts` and invoke each on the options struct.    
3. **Populate tools from registry** – for every init function in `globalToolsRegistry`:    
   * call it with the provided context & config;    
   * handle errors;    
   * skip nil results;    
   * if a whitelist is defined, filter by name;    
   * store the tool data in the map keyed by its definition name.    
4. **Finalize executor** – assign the built map to `toolsExecutor.Tools`, cache it globally, and return a pointer.  
  
### 4. Helper: `WithToolsWhitelist`    
```go  
func WithToolsWhitelist(tool ...string) ExecutorOption  
```  
Returns an option that appends supplied tool names to the whitelist inside an `ExecutorOptions`. This allows callers to specify which tools should be included when creating a new executor.  
  
---  
  
All parts together provide a flexible way to build a `tools.ToolsExecutor` from a registry of tool‑initializers, optionally filtered by a user‑supplied whitelist.  
  
# pkg/tools/webReader.go  
## Package / Component    
**Name:** `tools`    
  
### Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides the context type used in the tool call. |  
| `encoding/json` | Handles JSON unmarshalling of input arguments. |  
| `github.com/Swarmind/libagent/internal/tools` | Supplies the `ToolData` struct and a global registry for tools. |  
| `webreader "github.com/Swarmind/libagent/internal/tools/webReader"` | Provides the core logic to fetch and convert a URL into markdown text. |  
| `github.com/Swarmind/libagent/pkg/config` | Holds configuration flags (e.g., `WebReaderDisable`). |  
| `github.com/tmc/langchaingo/llms` | Supplies the LLM function definition type used for registering the tool. |  
  
---  
  
### External Data / Input Sources  
* **Configuration** – The tool is enabled/disabled via `cfg.WebReaderDisable`.    
* **Input JSON** – Expects a single field `"url"` that contains the URL to read.    
  
---  
  
### TODOs  
No explicit TODO comments are present in this file.  
  
---  
  
## Summary of Major Code Parts  
  
### 1. Function Definition (`WebReaderDefinition`)  
Defines an LLM‑compatible function named `webReader`.    
* **Name** – `"webReader"`.    
* **Description** – Explains that the tool reads a URL and returns markdown text.    
* **Parameters** – Expects an object with a single string property `"url"`.  
  
### 2. Argument Struct (`WebReaderArgs`)  
Simple struct used to unmarshal JSON input into a Go value:  
```go  
type WebReaderArgs struct {  
    URL string `json:"url"`  
}  
```  
The tag ensures the field is mapped from the key `"url"` in the incoming JSON.  
  
### 3. Tool Implementation (`WebReaderTool.Call`)  
Implements the actual work of the tool:  
1. Unmarshals the input JSON into a `WebReaderArgs` instance.  
2. Calls `webreader.ProcessUrl` with the provided URL and returns its result as a string (markdown text).    
The method signature matches the expected callback type for the registry.  
  
### 4. Registration (`init`)  
During package initialization, appends a closure to `globalToolsRegistry`.    
* The closure checks if the tool is disabled via `cfg.WebReaderDisable`.  
* It creates an instance of `WebReaderTool` and returns a populated `tools.ToolData` containing:  
  * `Definition`: the LLM function definition.  
  * `Call`: the method defined above.  
  
This makes the tool discoverable by the rest of the application.  
  
---  
  
All parts together provide a lightweight web‑reading utility that can be invoked through an LLM interface, configured via `config.Config`, and registered globally for use in higher‑level workflows.  
  
# pkg/tools/windows_executor.go  
# Package / Component    
**tools**  
  
## Imports  
```go  
import (  
	"context"  
	"encoding/json"  
	"fmt"  
	"os"  
	"strings"  
	"time"  
  
	"github.com/Swarmind/libagent/internal/tools"  
	"github.com/Swarmind/libagent/pkg/config"  
	"github.com/ThomasRooney/gexpect"  
	"github.com/google/uuid"  
	"github.com/rs/zerolog/log"  
  
	"github.com/tmc/langchaingo/llms"  
)  
```  
  
## External Data / Input Sources  
| Source | Description |  
|--------|-------------|  
| `config.Config` | The init function reads the following fields: <br>`CommandExecutorDisable` (bool) and `CommandExecutorCommands` (map[string]string). |  
  
## TODO List  
No explicit `TODO:` comments were found in this file.  
  
---  
  
# Summary of Major Code Parts  
  
### 1. Function Definition (`WCommandExecutorDefinition`)  
Defines a LangChainGo function that will be exposed to the LLM.    
- **Name**: `"windowsCommandExecutor"`    
- **Description**: Provides an interactive command‑execution tool for Windows CMD shells.    
- **Parameters**: Expects a JSON object with a single field `command` (string).    
  
### 2. Argument Struct (`WCommandExecutorArgs`)  
```go  
type WCommandExecutorArgs struct {  
	Command string `json:"command"`  
}  
```  
Used to unmarshal the input payload for the tool.  
  
### 3. Tool Implementation (`WCommandExecutorTool`)  
A stateful tool that keeps a temporary directory, an expect subprocess and a prompt string.  
  
#### Call Method  
```go  
func (s *WCommandExecutorTool) Call(ctx context.Context, input string) (string, error)  
```  
- Unmarshals the JSON payload into `commandExecutorArgs`.  
- Delegates to `RunCommand`.  
  
#### RunCommand Method  
Handles the actual execution:  
1. Creates a temporary directory if not already present.  
2. Spawns an interactive CMD session (`cmd /k`) in that directory.  
3. Sets up environment variables, a unique prompt (UUID), and sends the command string.  
4. Waits for the prompt to appear, collects output, parses it into lines, trims the trailing prompt, and returns the resulting string.  
  
#### Cleanup Method  
```go  
func (s *WCommandExecutorTool) Cleanup() error  
```  
- Removes the temporary directory if it exists.  
- Closes the subprocess.  
  
### 4. Registration (`init` function)  
Registers this tool in a global registry:  
- Checks `cfg.CommandExecutorDisable`; if false, creates an instance of `WCommandExecutorTool`.  
- Builds a command list string from `cfg.CommandExecutorCommands`, appending each command and its description.  
- Extends the definition’s description with that list.  
- Returns a `*tools.ToolData` containing the function definition, Call, and Cleanup callbacks.  
  
---  
  
All parts together provide an LLM‑exposed tool for executing arbitrary Windows CMD commands in a temporary session, with configuration driven command lists.  
  
