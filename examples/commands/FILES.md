# examples/commands/main.go  
## Package Summary: `main`  
  
This package demonstrates the usage of a command executor with whitelisted tools, specifically `ReWOOToolDefinition` and `CommandExecutorDefinition`. The main function initializes configurations, sets up logging, creates a tools executor instance, prepares a ReWOO query (defined in the `Prompt` constant), executes it using the tool executor, and prints the result.  
  
**Imports:**  
  
*   `context`: For managing context within asynchronous operations.  
*   `encoding/json`: For marshaling data into JSON format for tool execution.  
*   `fmt`: For formatted printing to standard output.  
*   `os`: For interacting with the operating system, such as setting up logging to stderr.  
*   `github.com/Swarmind/libagent/pkg/config`: For loading configuration settings.  
*   `github.com/Swarmind/libagent/pkg/tools`: For defining and executing whitelisted tools (ReWOO and CommandExecutor).  
*   `github.com/rs/zerolog`: For structured logging.  
*   `github.com/rs/zerolog/log`: For accessing the global logger instance.  
  
**External Data / Input Sources:**  
  
*   Configuration loaded via `config.NewConfig()`. The exact source of this configuration is not defined within this file but assumed to be external (e.g., environment variables, a config file).  
*   The `Prompt` constant defines a multi-step action plan that will be executed by the ReWOO tool. This prompt contains instructions for file manipulation, network checks, script execution, and git operations.  
  
**Major Code Parts:**  
  
1.  **Initialization & Logging**: Sets up global logging to stderr using zerolog with debug level.  
2.  **Configuration Loading**: Loads configuration settings using `config.NewConfig()`. Errors during loading are fatal.  
3.  **Tools Executor Setup**: Creates a tools executor instance, whitelisting only the ReWOO and CommandExecutor tools. The executor is deferred for cleanup to prevent resource leaks.  
4.  **ReWOO Query Preparation**: Defines a `Prompt` constant containing a series of commands to be executed by the ReWOO tool. This prompt includes file operations, network checks, script execution, git configuration, and commit actions. The query is marshaled into JSON format for passing to the tool executor.  
5.  **Tool Execution & Result Handling**: Calls the `ReWOOToolDefinition` with the prepared JSON query using the tools executor. Errors during tool invocation are fatal. The result of the execution (presumably a report as described in the prompt) is printed to standard output.  
  
**TODOs:**  
  
*   No explicit TODO comments found within this file. However, the reliance on external configuration and whitelisted tools implies potential future work related to managing these dependencies securely and efficiently.  
  
