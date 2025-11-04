# examples/generic/main.go  
# Package / Component    
**Name:** `main` (root package of a Go executable)    
  
## Imports  
| Import | Purpose |  
|--------|---------|  
| `context` | Provides the background context for agent execution |  
| `fmt` | Used to print the final result |  
| `os` | Accesses OS-level functions, e.g. stderr output |  
| `github.com/Swarmind/libagent/pkg/agent/generic` | Generic agent implementation |  
| `github.com/Swarmind/libagent/pkg/config` | Configuration handling (AI URL, token, model, etc.) |  
| `github.com/Swarmind/libagent/pkg/tools` | Tool executor and tool definitions |  
| `github.com/rs/zerolog` | Logging library |  
| `github.com/rs/zerolog/log` | Logger instance |  
| `github.com/tmc/langchaingo/llms/openai` | OpenAI LLM client |  
  
## External Data / Input Sources  
* **Configuration** – loaded via `config.NewConfig()`.    
  * Reads AI URL, token, model and default call options from a config file or environment.    
* **OpenAI API** – accessed through the client created by `openai.New(...)`.    
* **Tools** – four tool definitions are whitelisted: ReWOOTool, SemanticSearch, DDGSearch, WebReader.    
  
## TODOs  
No explicit `TODO:` comments were found in this file.  
  
---  
  
# Summary of Major Code Parts  
  
### Prompt Definition  
```go  
const Prompt = `Please use rewoo tool with the next prompt:  
Using semantic search tool, which can search across various code from the project  
collections find out the telegram library name in the code file contents for the project called "Hellper".  
Extract it from the given code and use a web search to find the pkg.go.dev documentation for it.  
Give me the URL for it.`  
```  
* A single string constant that will be fed into the agent’s LLM.  
  
### `main()` – Program Entry Point  
1. **Logging Setup**    
   * Sets global log level to debug and configures the logger to write to stderr via a console writer.  
  
2. **Configuration Load**    
   ```go  
   cfg, err := config.NewConfig()  
   ```  
   * Reads configuration values (AI URL, token, model, etc.) into `cfg`.  
  
3. **Context Creation**    
   ```go  
   ctx := context.Background()  
   ```  
  
4. **Agent Instantiation**    
   ```go  
   agent := generic.Agent{}  
   ```  
  
5. **LLM Client Setup**    
   * Creates an OpenAI client with base URL, token, model and API version from the config.  
   * Assigns this LLM to the agent (`agent.LLM = llm`).  
  
6. **Tools Executor Creation**    
   ```go  
   toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(  
       tools.ReWOOToolDefinition.Name,  
       tools.SemanticSearchDefinition.Name,  
       tools.DDGSearchDefinition.Name,  
       tools.WebReaderDefinition.Name,  
   ))  
   ```  
   * Builds a tool executor that knows which tools to use.  
   * The executor is attached to the agent (`agent.ToolsExecutor = toolsExecutor`).  
   * A deferred cleanup call ensures resources are released after execution.  
  
7. **Agent Run**    
   ```go  
   result, err := agent.SimpleRun(ctx,  
       Prompt, config.ConifgToCallOptions(cfg.DefaultCallOptions)...,  
   )  
   ```  
   * Executes the agent with the prompt and default call options.  
   * Prints the resulting output to stdout.  
  
---  
  
