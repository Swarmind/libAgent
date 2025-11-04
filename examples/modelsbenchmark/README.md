I’m ready to draft a concise, structured summary of your package once you share the files that make up the code base.  
Below is a template that follows the style shown in the example – just replace the placeholders with the real data from your project.

---

# Package / Component Overview  
**Package name:** `…`  

## Imports  
| Import | Purpose |
|--------|---------|
| `…` | … |
| `…` | … |
| `…` | … |

## External Data / Input Sources  
- **Configuration file(s)** – e.g. `config.yaml`, `settings.json`, etc.  
- **Command‑line flags** – e.g. `--verbose`, `--dry-run`.  
- **Environment variables** – e.g. `APP_ENV`, `LOG_LEVEL`.

## Project Package Structure (files & paths)  

```
<root>/
├── main.go
├── cmd/
│   └── run.go
├── pkg/
│   ├── config/
│   │   └── loader.go
│   └── tools/
│       └── executor.go
└── README.md
```

*(Adjust the tree to match your actual layout.)*

## Summary of Major Code Parts  

### 1. Global logger & configuration setup  
```go
zerolog.SetGlobalLevel(zerolog.InfoLevel)
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
cfg, err := config.NewConfig()
```
* Sets the global logging level to `Info`.  
* Directs zerolog output to standard error.  
* Loads a new configuration instance that will be reused for each operation.

### 2. Context creation  
```go
ctx := context.Background()
```
Creates a background context used by all tool calls.

### 3. Iteration over … (models, tasks, etc.)  
The `for …` loop performs the following per item:
1. **Logging** – prints which item is being processed.  
2. **Configuration update** – assigns the current value to `cfg`.  
3. **Tools executor creation** – builds a `toolsExecutor` with a whitelist containing two tool definitions (`…Definition`, `…Definition`).  
4. **Query preparation** – marshals the prompt into JSON bytes and calls the first tool.  
5. **Result handling** – prints the returned result, cleans up the executor, and sleeps for … if not on the last item.

### 4. Tool call details  
```go
toolArgs := tools.…ToolArgs{ Query: Prompt }
toolArgsBytes, err := json.Marshal(toolArgs)
result, err := toolsExecutor.CallTool(ctx,
    tools.…ToolDefinition.Name,
    string(toolArgsBytes),
)
```
* Builds a `…ToolArgs` struct with the prompt.  
* Serializes it to JSON for the tool call.  
* Executes the … tool and captures its output.

### 5. Cleanup & pacing  
```go
if err := toolsExecutor.Cleanup(); err != nil {
    log.Fatal().Err(err).Msg("tools executor cleanup")
}
…
time.Sleep(time.Minute * …)
```
Ensures that each item’s execution is cleaned up before the next iteration, and introduces a pause to allow the LocalAI watchdog (or similar) to unload the previous model.

---

This file orchestrates repeated calls to two tools for every item in `…`, using a single configuration object updated per iteration. The prompt outlines a series of actions that the first tool should perform, and the loop ensures each item is processed sequentially with appropriate logging and pacing.

---