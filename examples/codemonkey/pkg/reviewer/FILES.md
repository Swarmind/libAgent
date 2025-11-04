# examples/codemonkey/pkg/reviewer/reviewer.go  
# Package / Component    
**reviewer**  
  
The file implements a small helper that gathers information from a GitHub issue and returns it as a string. It is part of the *reviewer* package, which orchestrates calls to external tools (ReWOOTool, SemanticSearch, DDGSearch) via a `ToolsExecutor`.  
  
## Imports    
| Package | Purpose |  
|---------|---------|  
| `context` | Provides background context for tool execution |  
| `encoding/json` | Serialises the query struct before sending it to the executor |  
| `fmt` | Formats the prompt string in `CreatePrompt` |  
| `os` | Writes logs to standard error output |  
| `github.com/Swarmind/libagent/pkg/config` | Loads configuration for the executor |  
| `github.com/Swarmind/libagent/pkg/tools` | Tool definitions and executor helpers |  
| `github.com/rs/zerolog` | Logging level control |  
| `github.com/rs/zerolog/log` | Logger instance |  
  
## External data / input sources    
* **Configuration** – loaded via `config.NewConfig()`; this supplies the executor with runtime settings.    
* **Tool whitelist** – a slice of tool names (`ReWOOToolDefinition`, `DDGSearchDefinition`, `SemanticSearchDefinition`) that will be used by the executor.    
* **Issue & repo name** – passed into `GatherInfo` and forwarded to `CreatePrompt`.    
  
## TODOs    
No explicit `TODO:` comments are present in this file, but the following actions could be considered for future work:    
1. Add error handling for JSON marshaling failures.    
2. Cache or reuse the executor instead of creating a new one each call.  
  
## Summary of major code parts  
  
### `GatherInfo`  
* Sets global log level to debug and configures the logger.  
* Loads configuration, creates a context, and builds a whitelist of tools.  
* Instantiates a `ToolsExecutor` with the given context, config, and tool list.  
* Builds a query struct (`ReWOOToolArgs`) using `CreatePrompt`, marshals it into JSON, and calls the ReWOO tool via the executor.  
* Returns the raw string result from the tool call.  
  
### `CreatePrompt`  
* Constructs a detailed prompt for the ReWOO tool.    
  * The prompt instructs the AI to research an issue in a repository, using semantic search and DDG search tools.    
  * It includes placeholders for the issue title, repo name, and the names of the two tools that will be used.    
* Returns the formatted string ready for JSON serialization.  
  
The file is self‑contained: it pulls together configuration, tool execution, and prompt generation to produce a single string result that can be consumed by other components in the *reviewer* package.  
  
