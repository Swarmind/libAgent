# pkg/agent/generic/agent.go  
## Package: `generic`  
  
**Imports:**  
  
*   `context`: Standard Go context package for managing request lifecycles.  
*   `github.com/Swarmind/libagent/internal/tools`: Custom tools executor implementation.  
*   `github.com/tmc/langchaingo/llms`: LangChain LLM interface definitions.  
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI LLM integration for LangChain.  
  
**External Data / Inputs:**  
  
*   The `Agent` struct relies on an external `ToolsExecutor` instance (`tools.ToolsExecutor`) to manage and execute tools.  
*   LLM calls depend on the configuration of the injected `openai.LLM` instance (API key, model name, etc.).  
*   Input is provided via either a slice of `llms.MessageContent` for `Run()` or a string for `SimpleRun()`.  
  
**TODOs:**  
  
None found in this code snippet.  
  
### Agent Structure and Initialization  
  
The `Agent` struct encapsulates an OpenAI LLM (`openai.LLM`) and a tools executor (`tools.ToolsExecutor`). The `toolsList` field is lazily initialized to avoid unnecessary tool loading if not used.  Tool lists are retrieved from the `ToolsExecutor`.  
  
### Run Method  
  
The `Run()` method executes the core agent logic:  
  
1.  It retrieves or initializes the list of tools using the `ToolsExecutor`.  
2.  It calls the OpenAI LLM (`a.LLM.GenerateContent`) with provided state and tool options.  
3.  If the LLM response includes tool calls, it processes them via the `ToolsExecutor` to obtain updated content.  
4.  Finally, it returns the combined or processed content as a `llms.MessageContent`.  
  
### SimpleRun Method  
  
The `SimpleRun()` method provides a simplified interface for running the agent with a single string input:  
  
1.  It retrieves or initializes the list of tools using the `ToolsExecutor`.  
2.  It calls the OpenAI LLM (`a.LLM.GenerateContent`) with provided input as human message content and tool options.  
3.  If the LLM response includes tool calls, it processes them via the `ToolsExecutor` to obtain updated content.  
4.  Finally, it returns the combined or processed content as a string.  
  
