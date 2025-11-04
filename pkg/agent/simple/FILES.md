# pkg/agent/simple/agent.go  
# Package / Component    
**Name:** `simple`    
  
## Imports    
```go  
import (  
	"context"  
  
	"github.com/tmc/langchaingo/llms"  
	"github.com/tmc/langchaingo/llms/openai"  
)  
```  
* `context` – standard Go context package for request handling.    
* `github.com/tmc/langchaingo/llms` – core LangChain‑Go types and helpers.    
* `github.com/tmc/langchaingo/llms/openai` – OpenAI LLM implementation used by the agent.  
  
## External Data / Input Sources    
| Function | Parameters | Description |  
|----------|------------|-------------|  
| `Agent.Run` | `ctx context.Context`, `state []llms.MessageContent`, `opts ...llms.CallOption` | Takes a slice of message contents (the conversation state) and optional call options. |  
| `Agent.SimpleRun` | `ctx context.Context`, `input string`, `opts ...llms.CallOption` | Accepts a single human‑typed input string and optional call options. |  
  
## TODOs    
No explicit TODO comments are present in the file.  
  
---  
  
# Summary of Major Code Parts  
  
## 1. Agent struct  
```go  
type Agent struct {  
	LLM *openai.LLM  
}  
```  
* Holds a pointer to an `openai.LLM` instance, allowing the agent to call OpenAI’s LLM API directly.  
  
## 2. Run method  
```go  
func (a *Agent) Run(  
	ctx context.Context,  
	state []llms.MessageContent,  
	opts ...llms.CallOption,  
) (llms.MessageContent, error) {  
	response, err := a.LLM.GenerateContent(  
		ctx, state,  
	)  
	if err != nil {  
		return llms.MessageContent{}, err  
	}  
  
	content := response.Choices[0].Content  
  
	return llms.TextParts(llms.ChatMessageTypeAI, content), nil  
}  
```  
* Executes an LLM request with the provided conversation `state`.    
* Handles errors and returns the first choice’s content wrapped as a chat message of type AI.  
  
## 3. SimpleRun method  
```go  
func (a *Agent) SimpleRun(  
	ctx context.Context,  
	input string,  
	opts ...llms.CallOption,  
) (string, error) {  
	response, err := a.LLM.GenerateContent(ctx,  
		[]llms.MessageContent{  
			llms.TextParts(llms.ChatMessageTypeHuman,  
				input,  
			)},  
	)  
	if err != nil {  
		return "", err  
	}  
  
	content := response.Choices[0].Content  
  
	return content, nil  
}  
```  
* Provides a convenience wrapper that sends a single human‑typed input string to the LLM and returns the raw output string.  
  
---  
  
The file defines an `Agent` type that encapsulates an OpenAI LLM instance and offers two helper methods for generating chat responses: one that accepts a full conversation state (`Run`) and another that takes just a plain string input (`SimpleRun`).  
  
