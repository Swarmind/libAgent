# pkg/agent/generic/agent.go  
**Package name**    
`generic`  
  
---  
  
### Imports  
```go  
import (  
	"context"  
  
	"github.com/Swarmind/libagent/internal/tools"  
	"github.com/tmc/langchaingo/llms"  
	"github.com/tmc/langchaingo/llms/openai"  
)  
```  
* `context` – standard Go context package for request handling.    
* `github.com/Swarmind/libagent/internal/tools` – provides the `ToolsExecutor` type used to manage tool calls.    
* `github.com/tmc/langchaingo/llms` – core LLM interface and helper functions.    
* `github.com/tmc/langchaingo/llms/openai` – concrete OpenAI LLM implementation.  
  
---  
  
### External data / input sources  
| Source | Description |  
|--------|-------------|  
| `openai.LLM` | The underlying language‑model engine that generates content. |  
| `tools.ToolsExecutor` | Executes tool calls returned by the LLM. |  
| `[]llms.Tool` | List of tools to be passed to the LLM as options. |  
  
---  
  
### TODOs  
No explicit `TODO:` comments are present in this file, but a section is kept for future additions.  
  
---  
  
## Summary of major code parts  
  
#### 1. `Agent` struct    
```go  
type Agent struct {  
	LLM           *openai.LLM  
	ToolsExecutor *tools.ToolsExecutor  
	toolsList     *[]llms.Tool  
}  
```  
* Holds a reference to an OpenAI LLM, a tools executor, and a pointer to the list of tools that will be supplied to the LLM.  
  
#### 2. `Run` method    
```go  
func (a *Agent) Run(  
	ctx context.Context,  
	state []llms.MessageContent,  
	opts ...llms.CallOption,  
) (llms.MessageContent, error)  
```  
* Ensures `toolsList` is initialized and populated if empty.  
* Adds the tool list to the LLM options via `llms.WithTools`.  
* Calls `a.LLM.GenerateContent` with the current state and options.  
* Processes any tool calls returned by the LLM using `ToolsExecutor.ProcessToolCalls`.  
* Returns the first choice’s content as an AI chat message.  
  
#### 3. `SimpleRun` method    
```go  
func (a *Agent) SimpleRun(  
	ctx context.Context,  
	input string,  
	opts ...llms.CallOption,  
) (string, error)  
```  
* Similar to `Run`, but accepts a single input string instead of a slice of message content.  
* Wraps the input into a human‑type chat message before passing it to the LLM.  
* Returns the processed content directly as a plain string.  
  
---  
  
These two methods provide convenient ways to run an OpenAI model with or without pre‑built state, automatically handling tool execution and returning the final AI response.  
  
