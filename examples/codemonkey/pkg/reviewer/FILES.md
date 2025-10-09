# examples/codemonkey/pkg/reviewer/reviewer.go  
## Package: `reviewer` Summary  
  
This package provides functionality for gathering information about a GitHub issue using external tools and creating an actionable plan for developers. It leverages the `libagent` library for configuration and tool execution, along with logging via `zerolog`. The core function is `GatherInfo`, which orchestrates the process of researching an issue, synthesizing findings, and returning a structured result string.  
  
**Imports:**  
  
*   `context`: For managing context in asynchronous operations.  
*   `encoding/json`: For marshaling data to JSON format for tool communication.  
*   `fmt`: For formatted string output (prompt creation).  
*   `os`: For interacting with the operating system, specifically standard error stream logging.  
*   `github.com/Swarmind/libagent/pkg/config`: For loading configuration settings.  
*   `github.com/Swarmind/libagent/pkg/tools`: For executing external tools (ReWOOTool, DDGSearch, SemanticSearch).  
*   `github.com/rs/zerolog`: For structured logging.  
*   `github.com/rs/zerolog/log`: For global logger configuration and usage.  
  
**External Data / Input Sources:**  
  
*   `issue` (string): The GitHub issue string to analyze.  
*   `repoName` (string): The name of the repository containing the issue.  
*   Configuration loaded from `config.NewConfig()`.  The exact source of this configuration is not defined within this file but assumed to be external (e.g., environment variables, config files).  
  
**Major Code Parts:**  
  
### 1. `GatherInfo` Function: Core Logic  
  
This function orchestrates the entire process. It initializes logging, loads configuration, sets up a tools executor with whitelisted tools (`ReWOOTool`, `DDGSearch`, `SemanticSearch`), constructs a prompt based on the issue and repository name, calls the ReWOOTool with the prompt as input, and returns the result. Error handling is present throughout (fatal errors logged). The function also includes deferred cleanup of the tools executor to prevent resource leaks.  
  
### 2. `CreatePrompt` Function: Prompt Generation  
  
This function constructs a detailed prompt string that guides an AI agent in analyzing the GitHub issue. It uses formatted strings (`fmt.Sprintf`) to inject the issue and repository name into the prompt template. The prompt instructs the agent to research using specific tools (SemanticSearch, DDGSearch), synthesize findings, and produce a structured output including issue summary, desired outcome, relevant information, affected files, and code analysis snippets with comments.  The prompt explicitly warns against considering TODOs or commented-out code as instructions.  
  
**TODO Comments:**  
  
There are no explicit `// TODO:` style comments in the provided code snippet. However, the prompt itself contains implicit "TODO" like directives for the AI agent (e.g., "RESEARCH," "SYNTHESIS & ANALYSIS," "CREATE THE OUTPUT"). These aren't actionable items *within* the Go code but rather instructions for an external process consuming the generated prompt.  
  
