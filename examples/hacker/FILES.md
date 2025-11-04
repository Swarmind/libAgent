# examples/hacker/main.go  
**Package & Imports**    
- **Package name:** `main`    
- **Imports:**  
  - Standard library: `context`, `encoding/json`, `fmt`, `os`  
  - Project packages: `github.com/Swarmind/libagent/pkg/config`, `github.com/Swarmind/libagent/pkg/tools`  
  - Logging: `github.com/rs/zerolog`, `github.com/rs/zerolog/log`  
  
---  
  
**External Data Sources**    
- Environment variables:  
  - `HACKER_MSF_RHOST` – target address (RHOST)  
  - `HACKER_MSF_LHOST` – local host address (LHOST)  
  
These values are used to build the command string for the ReWOOTool.  
  
---  
  
**TODO List**    
| # | Comment |  
|---|----------|  
| 1 | `//func Act() {` – placeholder for a potential future wrapper function. |  
  
---  
  
## Code Overview  
  
### Logging & Configuration Setup  
```go  
zerolog.SetGlobalLevel(zerolog.DebugLevel)  
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  
```  
Initializes zerolog at debug level and configures the logger to write to standard error.  
  
### Config Loading  
```go  
cfg, err := config.NewConfig()  
if err != nil {  
    log.Fatal().Err(err).Msg("new config")  
}  
```  
Creates a new configuration object; errors are logged fatally if creation fails.  
  
### Tools Executor Creation  
```go  
ctx := context.Background()  
  
toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(  
    tools.ReWOOToolDefinition.Name,  
    tools.CommandExecutorDefinition.Name,  
    tools.NmapToolDefinition.Name,  
    tools.MsfSearchToolDefinition.Name,  
    tools.ExploitToolDefinition.Name,  
))  
```  
Builds a `ToolsExecutor` with a whitelist of five tool definitions (ReWOO, CommandExecutor, Nmap, MsfSearch, Exploit). The executor will later run the ReWOO command.  
  
### Environment Variable Retrieval  
```go  
rhost, exists := os.LookupEnv("HACKER_MSF_RHOST")  
if !exists {  
    log.Fatal().Msg("HACKER_MSF_RHOST env cannot be empty")  
}  
lhost, exists := os.LookupEnv("HACKER_MSF_LHOST")  
if !exists {  
    log.Fatal().Msg("HACKER_MSF_LHOST env cannot be empty")  
}  
  
fmt.Println(rhost, lhost)  
```  
Pulls the target and local host addresses from environment variables; prints them for debugging.  
  
### ReWOO Query Construction  
```go  
rewooQuery := tools.ReWOOToolArgs{  
    Query: fmt.Sprintf(Prompt, rhost, lhost),  
}  
```  
Creates a `ReWOOToolArgs` struct with a formatted query string that instructs the tool to scan ports and generate Metasploit search queries.  
  
### JSON Marshaling & Tool Execution  
```go  
rewooQueryBytes, err := json.Marshal(rewooQuery)  
if err != nil {  
    log.Fatal().Err(err).Msg("json marhsal rewooQuery")  
}  
  
result, err := toolsExecutor.CallTool(ctx,  
    tools.ReWOOToolDefinition.Name,  
    string(rewooQueryBytes),  
)  
```  
Serializes the query struct to JSON and passes it to the executor. The executor runs the ReWOO tool with the provided arguments.  
  
### Result Output  
```go  
fmt.Println(result)  
```  
Prints the raw result returned by the executor for quick inspection.  
  
---  
  
This file implements a command‑line helper that orchestrates several tools (ReWOO, Nmap, Metasploit search, etc.) to scan a host and generate exploit queries. It relies on environment variables for target data, builds a JSON payload, executes the tool chain, and logs the outcome.  
  
