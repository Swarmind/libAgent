# Package / Component  
**Name:** `simple` – an executable that demonstrates how to initialize and use a *simple* agent from the Swarmind library with an OpenAI LLM.

---

## Project File Structure  

```
examples/simple/
├─ main.go
├─ config/config.go
├─ agent/simple.go
└─ util/util.go
```

---

## Imports  
| Import | Purpose |
|--------|---------|
| `context` | Provides background context for the agent run. |
| `fmt` | Used for final output printing. |
| `os` | Access to OS-level functions (stderr writer). |
| `github.com/Swarmind/libagent/pkg/agent/simple` | The simple agent implementation that will be used. |
| `github.com/Swarmind/libagent/pkg/config` | Configuration handling – loads AI URL, token, model, etc. |
| `github.com/Swarmind/libagent/pkg/util` | Utility helpers (e.g., `RemoveThinkTag`). |
| `github.com/rs/zerolog` | Logging library for structured logs. |
| `github.com/rs/zerolog/log` | Convenience wrapper around zerolog. |
| `github.com/tmc/langchaingo/llms/openai` | OpenAI LLM client used by the agent. |

---

## External Data / Input Sources  
* **Configuration** – `config.NewConfig()` reads AI URL, token, model and default call options from a config file or environment.
* **Prompt** – constant string `"This is a test. Write OK in response."` that will be sent to the LLM.

---

## Environment Variables & Flags  
| Variable / Flag | Default / Example Value | Description |
|------------------|------------------------|-------------|
| `AIURL` | `https://api.openai.com/v1` | Base URL for the OpenAI API. |
| `AITOKEN` | `<your‑token>` | Authentication token for the LLM. |
| `MODEL` | `gpt-4o-mini` | Model name to use. |
| `DEFAULT_CALL_OPTIONS` | `[]string{"temperature=0.7","max_tokens=200"}` | Optional call options passed to the agent run. |

These can be overridden via command‑line flags (e.g., `--ai-url`, `--token`) or by setting environment variables before launching.

---

## Edge Cases for Launching  
| Scenario | Command |
|----------|---------|
| Run directly from source | `go run examples/simple/main.go` |
| Build a binary | `go build -o simple ./examples/simple` then `./simple` |
| Use a custom config file | `go run examples/simple/main.go --config=path/to/config.yaml` (if supported by the code) |

---

## Summary of Major Code Parts  

### 1. Logging Setup  
```go
zerolog.SetGlobalLevel(zerolog.DebugLevel)
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
```
Initializes zerolog at debug level and configures the logger to write to standard error.

### 2. Configuration Loading  
```go
cfg, err := config.NewConfig()
if err != nil {
    log.Fatal().Err(err).Msg("new config")
}
```
Creates a configuration object (`cfg`) that holds AI URL, token, model, and default call options. Errors are logged fatally.

### 3. Context & Agent Initialization  
```go
ctx := context.Background()
agent := simple.Agent{}
```
A background context is created for the agent run; an empty `simple.Agent` instance is prepared.

### 4. LLM Creation (OpenAI)  
```go
llm, err := openai.New(
    openai.WithBaseURL(cfg.AIURL),
    openai.WithToken(cfg.AIToken),
    openai.WithModel(cfg.Model),
    openai.WithAPIVersion("v1"),
)
if err != nil {
    log.Fatal().Err(err).Msg("new openai api llm")
}
agent.LLM = llm
```
Builds an OpenAI LLM client using the configuration values and attaches it to the agent.

### 5. Agent Run  
```go
result, err := agent.SimpleRun(ctx,
    Prompt, config.ConifgToCallOptions(cfg.DefaultCallOptions)...,
)
if err != nil {
    log.Fatal().Err(err).Msg("agent run")
}
```
Executes a simple run of the agent with the prompt and any default call options. Errors are logged.

### 6. Output Printing  
```go
fmt.Println(util.RemoveThinkTag(result))
```
Prints the result to standard output after stripping any “think” tags from the response.

---

## Notes on Code Relations  

* `main.go` orchestrates the whole flow: configuration → agent → LLM → run → print.
* The `config` package supplies a struct that is passed into the OpenAI client; its fields are used directly in the `openai.New()` call.
* The `agent/simple.go` file defines the `Agent` type and its method `SimpleRun`, which expects a context, prompt string, and optional call options.  
* Utility functions from `util/util.go` (e.g., `RemoveThinkTag`) are used to clean up the LLM response before printing.

---