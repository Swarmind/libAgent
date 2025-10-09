# examples/hacker/main.go  
## Package Summary: `main`  
  
This package demonstrates the usage of a command executor with whitelisted tools, specifically focusing on executing reconnaissance tasks using ReWOO and other security-related utilities (Nmap, Metasploit search). The primary function orchestrates tool execution based on a predefined prompt.  
  
**Imports:**  
  
*   `context`: For managing operation context.  
*   `encoding/json`: For serializing arguments to tools.  
*   `fmt`: For printing results.  
*   `os`: For standard output redirection (logging).  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration loading.  
*   `github.com/Swarmind/libagent/pkg/tools`: Tool execution framework.  
*   `github.com/rs/zerolog`: Structured logging.  
  
**External Data / Inputs:**  
  
*   The `Prompt` constant defines the reconnaissance query (scanning 172.86.95.138 for open ports and generating Metasploit search queries). This is hardcoded but could be externalized via configuration or command-line arguments in a real application.  
*   Configuration loaded from an unspecified source using `config.NewConfig()`. The exact config format isn't defined here, so it's assumed to provide necessary settings for the tools executor.  
  
**Whitelisted Tools:**  
  
The code explicitly whitelists these tools:  
  
*   `ReWOOToolDefinition.Name`: Executes a query (defined in `Prompt`) using ReWOO.  
*   `CommandExecutorDefinition.Name`: Generic command execution capability.  
*   `NmapToolDefinition.Name`: Network scanning with Nmap.  
*   `MsfSearchToolDefinition.Name`: Metasploit search queries based on discovered services.  
  
**Major Code Parts:**  
  
1.  **Initialization**: Sets up logging to standard error and loads configuration using `config.NewConfig()`.  
2.  **Tools Executor Setup**: Creates a `toolsExecutor` instance with the specified whitelist of tools, ensuring only approved commands can be executed. The executor is deferred for cleanup.  
3.  **ReWOO Execution**: Constructs a ReWOO query from the hardcoded `Prompt`, serializes it to JSON, and calls the tool using `toolsExecutor.CallTool()`.  
4.  **Output**: Prints the raw result returned by ReWOO to standard output.  
  
**TODOs:**  
  
There are no explicit TODO comments in this code snippet. However, potential improvements include:  
  
*   Externalizing the `Prompt` for dynamic configuration.  
*   Adding error handling beyond fatal logging (e.g., retries or fallback mechanisms).  
*   Implementing input validation to prevent malicious queries from being executed via ReWOO.  
  
