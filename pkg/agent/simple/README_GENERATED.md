## Package: `simple` Summary

This package implements a basic OpenAI agent using LangChainGo for running prompts and generating content. It provides two methods, `Run` (for message-based interactions) and `SimpleRun` (for direct string input), both leveraging an injected OpenAI LLM instance to perform the actual generation. The core logic revolves around formatting inputs into appropriate structures for the LLM and handling potential errors during content creation.

**Imports:**

*   `context`: For managing request lifecycles.
*   `github.com/tmc/langchaingo/llms`: LangChainGo's base LLM interface and message structures.
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI-specific LLM implementation within LangChainGo.

**Configuration:**

*   Requires an initialized `openai.LLM` instance to be injected into the `Agent`. This LLM must be properly configured with API keys and other necessary settings for accessing the OpenAI service.
*   The `Run` method accepts a slice of `llms.MessageContent`, representing a conversation history or prompt context.
*   The `SimpleRun` method takes a single string as input, which is treated as a direct user query.

**File Structure:**

```
pkg/agent/simple/
├── agent.go
```

**Code Entities and Relations:**

*   `Agent`: Struct containing an OpenAI LLM instance (`LLM *openai.LLM`). This dependency injection pattern allows for easy swapping of different LLM implementations if needed.
*   `Run(ctx context.Context, messages []llms.MessageContent) (string, error)`: Handles message-based interactions by passing a slice of `llms.MessageContent` to the underlying LLM's `GenerateContent` method. The response is then extracted and returned as an AI-generated message. Error handling is included for LLM failures.
*   `SimpleRun(ctx context.Context, prompt string) (string, error)`: Provides a simplified interface for direct string input, converting it into a single human-role message before sending it to the LLM. Similar error handling applies.