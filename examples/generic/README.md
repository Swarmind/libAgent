## Package Summary: `main`

This package demonstrates initializing and running a generic agent using the `libagent` library to find Telegram library documentation via semantic search, web search (DuckDuckGo), and OpenAI's LLM. It loads configuration from environment variables or config files for OpenAI API access.

**Imports:**

*   `context`: For request context management.
*   `fmt`: Formatted printing.
*   `os`: OS interaction (stderr output).
*   `github.com/Swarmind/libagent/pkg/agent/generic`: Core agent functionality.
*   `github.com/Swarmind/libagent/pkg/config`: Configuration loading.
*   `github.com/Swarmind/libagent/pkg/tools`: Tool definitions and execution.
*   `github.com/rs/zerolog`, `github.com/rs/zerolog/log`: Structured logging.
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI LLM integration.

**Configuration:**

The application relies on the following environment variables:

*   `AIURL`: Base URL for OpenAI API (e.g., `https://api.openai.com/v1`).
*   `AIToken`: OpenAI API authentication token.
*   `Model`: The name of the OpenAI model to use (e.g., "gpt-3.5-turbo").
*   `DefaultCallOptions`: Default options for LLM calls (not fully specified in snippet).

**Tools:**

The agent uses these tools:

*   `ReWOOToolDefinition`: Purpose unclear, likely text processing/rewriting.
*   `SemanticSearchDefinition`: Searches project code files.
*   `DDGSearchDefinition`: DuckDuckGo search tool.
*   `WebReaderDefinition`: Reads content from web pages.

**Execution Flow:**

1.  Initializes logging and loads configuration.
2.  Sets up OpenAI LLM with API credentials. Exits if initialization fails.
3.  Creates a `toolsExecutor` with allowed tools (ReWOO, Semantic Search, DDG Search, Web Reader).
4.  Runs the agent with a predefined prompt to find Telegram library documentation using semantic search and web search.
5.  Prints the final result to stdout.

**Project Structure:**

*   `main.go`: Contains all logic for initializing, configuring, running the agent, and printing results.