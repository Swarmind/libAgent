## Package Summary: `main`

This package defines a command-line application that executes reconnaissance tasks using whitelisted tools, including ReWOO, Nmap, and Metasploit search. The primary function orchestrates tool execution based on a predefined prompt.

**Imports:**

*   `context`: For managing operation context.
*   `encoding/json`: For serializing arguments to tools.
*   `fmt`: For printing results.
*   `os`: For standard output redirection (logging).
*   `github.com/Swarmind/libagent/pkg/config`: Configuration loading.
*   `github.com/Swarmind/libagent/pkg/tools`: Tool execution framework.
*   `github.com/rs/zerolog`: Structured logging.

**Environment Variables / Flags:**

The package relies on configuration loaded via `config.NewConfig()`, which may read from environment variables or command-line flags (not explicitly defined in this snippet). The exact configuration mechanism is unspecified, but it's assumed to provide necessary settings for the tools executor.

**File Structure:**

*   `main.go`: Contains the main application logic and tool execution orchestration.

**Major Code Parts:**

1.  **Initialization**: Sets up logging to standard error and loads configuration using `config.NewConfig()`.
2.  **Tools Executor Setup**: Creates a `toolsExecutor` instance with the specified whitelist of tools, ensuring only approved commands can be executed. The executor is deferred for cleanup.
3.  **ReWOO Execution**: Constructs a ReWOO query from the hardcoded `Prompt`, serializes it to JSON, and calls the tool using `toolsExecutor.CallTool()`.
4.  **Output**: Prints the raw result returned by ReWOO to standard output.

**Edge Cases:**

*   The application relies on external tools (ReWOO, Nmap, Metasploit) being installed and accessible in the system's PATH. Failure to meet this requirement will cause execution errors.
*   The hardcoded `Prompt` limits the reconnaissance scope. Modifying it requires code changes or an external configuration mechanism.

**Potential Issues:**

*   The lack of input validation could allow malicious queries if the `Prompt` is dynamically sourced from untrusted inputs.
*   Error handling is minimal (fatal logging). More robust error management might be needed for production use cases.