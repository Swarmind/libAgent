# examples/modelsbenchmark/main.go  
# Package / Component Overview    
**Package name:** `main`    
  
## Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides background context for tool execution |  
| `encoding/json` | Marshal query arguments to JSON |  
| `fmt` | Print results to stdout |  
| `os` | Access stderr for zerolog output |  
| `time` | Sleep between model calls |  
| `github.com/Swarmind/libagent/pkg/config` | Load and store configuration data |  
| `github.com/Swarmind/libagent/pkg/tools` | Create tool executor, whitelist tools, call tools |  
| `github.com/rs/zerolog` | Logging library |  
| `github.com/rs/zerolog/log` | Logger instance |  
  
---  
  
## External Data / Input Sources  
- **Model list** – a hard‑coded slice of model identifiers (`ModelList`) that will be iterated over.  
- **Prompt string** – a multiline plan used as the query for the ReWOOTool.    
- **Configuration** – loaded via `config.NewConfig()` and updated per iteration with the current model.  
  
---  
  
## TODO Comments  
No explicit `TODO:` markers are present in this file, but the comment at the top indicates that the example demonstrates multiple tool calls for different config values.  
  
---  
  
# Summary of Major Code Parts  
  
### 1. Global logger & configuration setup    
```go  
zerolog.SetGlobalLevel(zerolog.InfoLevel)  
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  
cfg, err := config.NewConfig()  
```  
* Sets the global logging level to `Info`.    
* Directs zerolog output to standard error.    
* Loads a new configuration instance that will be reused for each model.  
  
### 2. Context creation    
```go  
ctx := context.Background()  
```  
Creates a background context used by all tool calls.  
  
### 3. Iteration over models    
The `for idx, model := range ModelList` loop performs the following per model:  
1. **Logging** – prints which model is being processed.    
2. **Configuration update** – assigns the current model to `cfg.Model`.    
3. **Tools executor creation** – builds a `toolsExecutor` with a whitelist containing two tool definitions (`ReWOOToolDefinition`, `CommandExecutorDefinition`).    
4. **Query preparation** – marshals the prompt into JSON bytes and calls the ReWOO tool.    
5. **Result handling** – prints the returned result, cleans up the executor, and sleeps for 2 minutes if not on the last model.  
  
### 4. Tool call details    
```go  
rewooQuery := tools.ReWOOToolArgs{ Query: Prompt }  
rewooQueryBytes, err := json.Marshal(rewooQuery)  
result, err := toolsExecutor.CallTool(ctx,  
    tools.ReWOOToolDefinition.Name,  
    string(rewooQueryBytes),  
)  
```  
* Builds a `ReWOOToolArgs` struct with the prompt.    
* Serializes it to JSON for the tool call.    
* Executes the ReWOO tool and captures its output.  
  
### 5. Cleanup & pacing    
```go  
if err := toolsExecutor.Cleanup(); err != nil {  
    log.Fatal().Err(err).Msg("tools executor cleanup")  
}  
...  
time.Sleep(time.Minute * 2)  
```  
Ensures that each model’s execution is cleaned up before the next iteration, and introduces a 2‑minute pause to allow the LocalAI watchdog to unload the previous model.  
  
---  
  
This file orchestrates repeated calls to two tools for every model in `ModelList`, using a single configuration object updated per iteration. The prompt outlines a series of actions that the ReWOO tool should perform, and the loop ensures each model is processed sequentially with appropriate logging and pacing.  
  
