# examples/gitHelperLegacy/main.go  
## Package `main` Summary  
  
This package implements an example workflow using the ReWOO agent to manage GitHub issues, inspired by similar examples in the `libagent` library. The primary goal is to automate issue resolution for a given repository and issue description.  The code sets up tools like `ReWOOTool`, `CommandExecutor`, `SemanticSearchTool`, and `WebReaderTool` within a configurable environment.  
  
### Imports:  
  
*   **context:** For managing asynchronous operations.  
*   **encoding/json:** For serializing ReWOO query arguments.  
*   **flag:** For parsing command-line arguments (repo URL, issue input).  
*   **fmt:** For formatted output and string manipulation.  
*   **os:** For interacting with the operating system (stderr for logging).  
*   **path/filepath:** For extracting repository name from URL.  
*   **strings:** For manipulating strings (trimming repo URLs).  
*   **github.com/Swarmind/libagent/pkg/config:** For loading application configuration.  
*   **github.com/Swarmind/libagent/pkg/tools:** For defining and executing tools within the ReWOO agent workflow.  
*   **github.com/rs/zerolog, github.com/rs/zerolog/log:** For structured logging with debug-level output to stderr.  
  
### External Data & Inputs:  
  
*   **Repository URL (`-repo` flag):**  Defaults to `https://github.com/JackBekket/Reflexia`.  
*   **Issue Input (`-issue` flag):** Required; can be an issue number, full URL, or descriptive text.  
*   **Additional Context (`-additionalContext` flag):** Optional context for the agent's task.  
  
### TODOs:  
  
No explicit `TODO` comments are present in this code snippet. However, operational reminders at the end suggest prerequisites like installing and configuring `git` and GitHub CLI (`gh`) with authentication.  The SemanticSearchTool requires a pre-populated pgvector database collection named after the repository.  
  
### Major Code Parts Summary:  
  
1.  **Configuration & Logging:** Initializes application configuration using `config.NewConfig()` and sets up zerolog logging to stderr at debug level.  
2.  **Command-Line Argument Parsing:** Parses command-line flags for repo URL, issue input, and additional context. Validates that the issue input is provided.  
3.  **Tools Executor Setup:** Creates a `toolsExecutor` instance with whitelisted tools (ReWOO, CommandExecutor, SemanticSearchTool, WebReaderTool). Handles cleanup of the executor using defer.  
4.  **Repository Details Extraction:** Derives repository name and owner/repo string from the provided URL. Defaults to "Reflexia" if the default repo is used.  
5.  **Task Description Construction:** Constructs a detailed task description for the ReWOO agent, including instructions on understanding the issue, setting up the development environment (cloning, changing directories), finding and modifying files using `CommandExecutor` and `SemanticSearchTool`, verifying changes with `git status`, and handling multi-line modifications. Includes warnings about exceeding 4000 character limits for task descriptions.  
6.  **ReWOO Agent Execution:** Marshals the ReWOO query arguments into JSON, calls the ReWOO tool using the tools executor, prints the result to stdout, and provides operational reminders (git/gh installation, authentication, pgvector database setup).  
  
The code is designed as a self-contained example of how to integrate the ReWOO agent with external tools for automating GitHub issue resolution. The task description is highly detailed, providing step-by-step instructions for the agent to follow.  Error handling and logging are implemented throughout the process.  
  
