# pkg/agent/simple/agent.go  
## Package: `simple` Summary  
  
This package provides a simple agent interface for interacting with OpenAI LLMs via the LangChainGo framework. It offers two primary methods for running prompts: `Run` (for message-based interactions) and `SimpleRun` (for direct string input). The core functionality revolves around generating content using an injected OpenAI LLM instance.  
  
**Imports:**  
  
*   `context`: Standard Go context package for managing request lifecycles.  
*   `github.com/tmc/langchaingo/llms`: LangChainGo's base LLM interface and message structures.  
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI-specific LLM implementation within LangChainGo.  
  
**External Data / Inputs:**  
  
*   Requires an initialized `openai.LLM` instance to be injected into the `Agent`. This LLM must be properly configured with API keys and other necessary settings for accessing the OpenAI service.  
*   The `Run` method accepts a slice of `llms.MessageContent`, representing a conversation history or prompt context.  
*   The `SimpleRun` method takes a single string as input, which is treated as a direct user query.  
  
**Major Code Parts:**  
  
1.  **Agent Struct**: Defines the agent with an embedded OpenAI LLM instance (`LLM *openai.LLM`). This dependency injection pattern allows for easy swapping of different LLM implementations if needed.  
2.  **Run Method**: Handles message-based interactions by passing a slice of `llms.MessageContent` to the underlying LLM's `GenerateContent` method. The response is then extracted and returned as an AI-generated message. Error handling is included for LLM failures.  
3.  **SimpleRun Method**: Provides a simplified interface for direct string input, converting it into a single human-role message before sending it to the LLM. Similar error handling applies.  
  
**TODOs:**  
  
*   No TODO comments found in this file.  
  
