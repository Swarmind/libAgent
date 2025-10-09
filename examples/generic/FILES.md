# examples/generic/main.go  
## Package Summary: `main`  
  
This package demonstrates the initialization and usage of a generic tools-enabled agent using the `libagent` library. It leverages OpenAI's LLM for processing prompts and executes predefined tools to achieve a specific task (finding Telegram library documentation).  
  
**Imports:**  
  
*   `context`: For managing request contexts.  
*   `fmt`: For formatted printing.  
*   `os`: For interacting with the operating system, specifically standard error output.  
*   `github.com/Swarmind/libagent/pkg/agent/generic`: Core agent functionality.  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration loading and handling.  
*   `github.com/Swarmind/libagent/pkg/tools`: Tool definitions and execution.  
*   `github.com/rs/zerolog`, `github.com/rs/zerolog/log`: Structured logging.  
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI LLM integration.  
  
**External Data / Input Sources:**  
  
*   Configuration loaded from environment variables or a configuration file (using `config.NewConfig()`).  Specifically, the following config values are used:  
    *   `AIURL`: Base URL for the OpenAI API.  
    *   `AIToken`: Authentication token for the OpenAI API.  
    *   `Model`: The name of the OpenAI model to use (e.g., "gpt-3.5-turbo").  
    *   `DefaultCallOptions`: Default options passed to LLM calls.  
  
**Tools Used:**  
  
The code utilizes a predefined set of tools:  
  
*   `ReWOOToolDefinition`:  (Purpose unspecified in the provided snippet, likely related to rewriting or processing text).  
*   `SemanticSearchDefinition`: For searching across project collections (code files) for specific information.  
*   `DDGSearchDefinition`: DuckDuckGo search tool.  
*   `WebReaderDefinition`: Tool for reading content from web pages.  
  
**TODOs:**  
  
There are no explicit `TODO` comments in the provided code snippet.  
  
**Major Code Parts Summary:**  
  
1.  **Initialization & Configuration:** The program initializes logging, loads configuration (including OpenAI API credentials), and sets up a context.  
2.  **LLM Setup:** An OpenAI LLM instance is created using the loaded configuration parameters. Error handling ensures that if the LLM cannot be initialized, the program exits.  
3.  **Tools Executor Setup:** A `toolsExecutor` is instantiated with a whitelist of allowed tools (ReWOO, Semantic Search, DDG Search, Web Reader). The executor also includes deferred cleanup logic to release resources.  
4.  **Agent Execution:** The core functionality involves running the agent with a predefined prompt (`Prompt`). This prompt instructs the agent to use semantic search to find the Telegram library name in project code, then use web search (via DDG) to locate its pkg.go.dev documentation URL.  
5.  **Output:** The final result of the agent's execution is printed to standard output.  
  
The primary purpose of this file appears to be a demonstration or example of how to integrate and utilize tools within the `libagent` framework, specifically for code-related tasks involving OpenAI LLMs. It showcases a workflow where an LLM orchestrates tool usage based on a given prompt.  
  
