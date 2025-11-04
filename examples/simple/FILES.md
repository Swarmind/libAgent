# examples/simple/main.go  
# Package / Component    
**Name:** `main` – a simple executable that demonstrates how to initialize and use a *simple* agent from the Swarmind library with an OpenAI LLM.  
  
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
  
## TODOs  
No explicit TODO comments are present; the section can be expanded later if needed.  
  
---  
  
# Summary of Major Code Parts  
  
## 1. Logging Setup    
```go  
zerolog.SetGlobalLevel(zerolog.DebugLevel)  
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  
```  
Initializes zerolog at debug level and configures the logger to write to standard error.  
  
## 2. Configuration Loading    
```go  
cfg, err := config.NewConfig()  
if err != nil {  
    log.Fatal().Err(err).Msg("new config")  
}  
```  
Creates a configuration object (`cfg`) that holds AI URL, token, model, and default call options. Errors are logged fatally.  
  
## 3. Context & Agent Initialization    
```go  
ctx := context.Background()  
agent := simple.Agent{}  
```  
A background context is created for the agent run; an empty `simple.Agent` instance is prepared.  
  
## 4. LLM Creation (OpenAI)    
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
  
## 5. Agent Run    
```go  
result, err := agent.SimpleRun(ctx,  
    Prompt, config.ConifgToCallOptions(cfg.DefaultCallOptions)...,  
)  
if err != nil {  
    log.Fatal().Err(err).Msg("agent run")  
}  
```  
Executes a simple run of the agent with the prompt and any default call options. Errors are logged.  
  
## 6. Output Printing    
```go  
fmt.Println(util.RemoveThinkTag(result))  
```  
Prints the result to standard output after stripping any “think” tags from the response.  
  
---  
  
