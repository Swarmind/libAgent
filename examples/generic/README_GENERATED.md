**generic**

This package is an example demonstrating the use of the generic agent, which interacts with OpenAI LLM, handles configuration, and executes tools.

External Data/Config Options:
- `AIURL`: OpenAI API base URL (from config).
- `AIToken`: OpenAI API token (from config).
- `Model`: LLM model identifier (from config).
- `DefaultCallOptions`: Default call options for the agent (from config).
- `Tools` are restricted to: `ReWOOTool`, `SemanticSearch`, `DDGSearch`, `WebReader`.

TODOs/Comments:
- `config.ConifgToCallOptions` (likely a typo: should be `ConfigToCallOptions`).
- No error handling for `toolsExecutor.Cleanup()` (though `defer` ensures it runs).

Edge Cases:
- If `config.NewConfig()` fails, the program exits with a fatal log.
- If LLM initialization fails, the program exits with a fatal log.
- If tools executor initialization fails, the program exits with a fatal log.
- If `agent.SimpleRun()` fails, the program exits with a fatal log.