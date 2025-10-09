# examples/codemonkey/pkg/planner/planner.go  
## Package: `planner` Summary  
  
**Package Name:** `planner`  
  
**Imports:**  
  
*   `context`: For managing request contexts.  
*   `encoding/json`: For JSON serialization and deserialization.  
*   `os`: For interacting with the operating system (e.g., standard error).  
*   `github.com/Swarmind/libagent/pkg/config`: For configuration management.  
*   `github.com/Swarmind/libagent/pkg/tools`: For tool execution and whitelisting.  
*   `github.com/rs/zerolog`: For structured logging.  
*   `github.com/rs/zerolog/log`: For global logger access.  
  
**External Data & Inputs:**  
  
*   The package relies on a configuration object (`config.Config`) loaded from an unspecified source (likely environment variables or files).  
*   It takes string inputs for `PlanGitHelper` and `PlanCLIExecutor` functions, which are used as prompts to external tools via the `tools.ToolsExecutor`. The exact nature of these strings determines the output.  
  
**TODOs:**  
  
No explicit TODO comments found in this code snippet.  
  
### Function: `PlanGitHelper`  
  
This function orchestrates a process for transforming a "Reviewer Result" (a string input) into an executable action plan using an external tool (`tools.ReWOOToolDefinition`). It constructs a prompt based on the `lePromptGithelper` constant and the provided review text, serializes it as JSON, calls the tool with this payload, and returns the result. Error handling includes fatal logging if any step fails (config loading, tool execution, or empty results). The function also sets up global zerolog logging to stderr for debugging purposes.  
  
### Function: `PlanCLIExecutor`  
  
Similar to `PlanGitHelper`, this function takes a string input (`task`) and transforms it into executable CLI commands using an external tool (`tools.ReWOOToolDefinition`). It uses the `lePromptCLI` constant as part of the prompt, serializes it as JSON, calls the tool, and returns the result. Error handling is identical to `PlanGitHelper`, with fatal logging on failure. The function also configures zerolog logging to stderr for debugging.  
  
**Key Observations:**  
  
*   Both functions heavily rely on an external tool (`tools.ReWOOToolDefinition`) for processing prompts into executable instructions.  
*   The package uses a consistent pattern of loading configuration, constructing prompts, serializing them as JSON, calling the tool, and returning the result with fatal error handling.  
*   Logging is configured globally using zerolog to stderr.  
  
