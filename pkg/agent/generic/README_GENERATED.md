## Package: `generic`

This package implements an OpenAI-powered agent capable of executing tools based on LLM responses. The core functionality revolves around interacting with an external `ToolsExecutor` and an OpenAI LLM instance (`openai.LLM`).  The agent handles both structured message inputs (via `Run()`) and simple string prompts (via `SimpleRun()`), dynamically invoking tools as directed by the LLM's output.

**Configuration:**

*   **OpenAI API Key/Model:** The behavior is heavily dependent on the configured OpenAI LLM (`openai.LLM`).  API keys, model names, temperature settings, and other parameters must be set correctly for proper operation.
*   **ToolsExecutor:** The `tools.ToolsExecutor` instance provides access to available tools. Its configuration (e.g., tool definitions) dictates which actions the agent can perform.

**Files:**

*   `pkg/agent/generic/agent.go`: Contains the main `Agent` struct and its methods (`Run`, `SimpleRun`).

**Workflow:**

1.  The agent receives input (either a message slice or a string).
2.  It retrieves a list of available tools from the `ToolsExecutor`.
3.  It prompts the OpenAI LLM with the input, including tool descriptions in the prompt.
4.  If the LLM's response indicates tool usage, the agent extracts those calls and executes them via the `ToolsExecutor`.
5.  The results of tool execution are fed back into the LLM for further processing (if necessary).
6.  The final output is returned as either a message slice or a string.

**Edge Cases:**

*   If the OpenAI API key is invalid, requests will fail.
*   If the `ToolsExecutor` fails to load tools correctly, the agent may not be able to perform any actions.
*   LLM hallucinations can lead to incorrect tool calls or nonsensical responses.  Proper prompt engineering and safety measures are crucial.

**Potential Issues:**

The code relies heavily on external dependencies (OpenAI API, `ToolsExecutor`). Any instability in these components will directly impact the agent's behavior. The lazy initialization of tools might introduce minor performance overhead if tools are frequently accessed but not always used.