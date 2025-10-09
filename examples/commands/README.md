## Package Summary: `main`

This package defines a command-line application that executes a predefined sequence of actions using a tool executor with whitelisted tools (ReWOO and CommandExecutor). The main function initializes logging, loads configuration, sets up the tool executor, prepares a ReWOO query, executes it, and prints the result.  The core logic revolves around executing a hardcoded `Prompt` through the ReWOO tool.

**Imports:**

*   `context`: For managing context within asynchronous operations.
*   `encoding/json`: For marshaling data into JSON format for tool execution.
*   `fmt`: For formatted printing to standard output.
*   `os`: For interacting with the operating system, such as setting up logging to stderr.
*   `github.com/Swarmind/libagent/pkg/config`: For loading configuration settings.
*   `github.com/Swarmind/libagent/pkg/tools`: For defining and executing whitelisted tools (ReWOO and CommandExecutor).
*   `github.com/rs/zerolog`: For structured logging.
*   `github.com/rs/zerolog/log`: For accessing the global logger instance.

**Environment Variables / Configuration:**

The application relies on external configuration loaded via `config.NewConfig()`. The exact source of this configuration (e.g., environment variables, config file) is not defined within this file but assumed to be provided externally.  No specific environment variable names are hardcoded in the code.

**Command-Line Arguments:**

The application does not accept any command-line arguments directly. All behavior is determined by the hardcoded `Prompt` and external configuration loaded at runtime.

**File Structure:**

*   `main.go`: Contains the main entry point, initialization logic, tool executor setup, ReWOO query preparation, execution, and result printing.

**Code Logic Summary:**

1.  **Initialization**: Sets up global logging to stderr with debug level using zerolog.
2.  **Configuration Loading**: Loads configuration settings using `config.NewConfig()`. Errors during loading are fatal.
3.  **Tools Executor Setup**: Creates a tools executor instance, whitelisting only the ReWOO and CommandExecutor tools. The executor is deferred for cleanup to prevent resource leaks.
4.  **ReWOO Query Preparation**: Defines a `Prompt` constant containing a series of commands to be executed by the ReWOO tool. This prompt includes file operations, network checks, script execution, git configuration, and commit actions. The query is marshaled into JSON format for passing to the tool executor.
5.  **Tool Execution & Result Handling**: Calls the `ReWOOToolDefinition` with the prepared JSON query using the tools executor. Errors during tool invocation are fatal. The result of the execution (presumably a report as described in the prompt) is printed to standard output.