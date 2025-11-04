# examples/hacker/builder/builder.go  
# Package / Component Summary  
  
**Package name:** `main`    
The file implements a small command‑line helper that orchestrates a *ReWOO* tool call to generate, build, and run a Go binary.    
  
## Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides the background context for tool execution. |  
| `encoding/json` | Marshals the ReWOO query into JSON. |  
| `fmt` | Prints the final result to stdout. |  
| `os` | Accesses standard error output for logging. |  
| `github.com/Swarmind/libagent/pkg/config` | Loads configuration data for the executor. |  
| `github.com/Swarmind/libagent/pkg/tools` | Supplies tool definitions and execution logic. |  
| `github.com/rs/zerolog` | Logging library (global level). |  
| `github.com/rs/zerolog/log` | Convenience logger wrapper. |  
  
## External Data / Input Sources  
- **Configuration** – loaded via `config.NewConfig()`.    
- **ReWOO Tool Definition** – referenced by `tools.ReWOOToolDefinition.Name`.    
- **Command Executor Definition** – referenced by `tools.CommandExecutorDefinition.Name`.    
- **Builder Prompt** – a constant string that describes the step‑by‑step actions for the ReWOO tool.    
  
## TODOs  
| Line | Comment |  
|------|---------|  
| 1 | “Do not cleanup to use the final result” – currently the cleanup defer is commented out; consider adding it later. |  
  
## Main Function Flow (Markdown Subheaders)  
  
### `func main()`  
- **Logging Setup**    
  ```go  
  zerolog.SetGlobalLevel(zerolog.DebugLevel)  
  log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  
  ```  
  Sets the global logging level to *Debug* and directs logs to standard error.  
  
- **Configuration Loading**    
  ```go  
  cfg, err := config.NewConfig()  
  if err != nil { log.Fatal().Err(err).Msg("new config") }  
  ```  
  Creates a new configuration object; fatal on failure.  
  
- **Context Creation**    
  ```go  
  ctx := context.Background()  
  ```  
  
- **Tools Executor Instantiation**    
  ```go  
  toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg,  
      tools.WithToolsWhitelist(  
          tools.ReWOOToolDefinition.Name,  
          tools.CommandExecutorDefinition.Name,  
      ),  
  )  
  if err != nil { log.Fatal().Err(err).Msg("new tools executor") }  
  ```  
  Builds a `ToolsExecutor` that will run the ReWOO and Command Executor tools.  
  
- **ReWOO Query Preparation**    
  ```go  
  rewooQuery := tools.ReWOOToolArgs{  
      Query: BuilderPrompt,  
  }  
  rewooQueryBytes, err := json.Marshal(rewooQuery)  
  if err != nil { log.Fatal().Err(err).Msg("json marhsal rewooQuery") }  
  ```  
  
- **Tool Execution**    
  ```go  
  result, err := toolsExecutor.CallTool(ctx,  
      tools.ReWOOToolDefinition.Name,  
      string(rewooQueryBytes),  
  )  
  if err != nil { log.Fatal().Err(err).Msg("rewoo tool call") }  
  ```  
  Calls the ReWOO tool with the prepared JSON query.  
  
- **Result Output**    
  ```go  
  fmt.Println(result)  
  ```  
  
The code is intentionally minimal; it prints the raw result of the ReWOO tool to standard output. Future enhancements could include handling cleanup, error handling improvements, and integration with the Command Executor tool for building the binary.  
  
