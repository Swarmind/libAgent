```markdown
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

**Environment Variables/Configuration:**

The package depends on a loaded configuration (`config.Config`).  Specific environment variables or file paths for loading this config are not defined in the provided code snippet, but it's assumed to be handled elsewhere (e.g., via flags or default locations). The `tools.ToolsExecutor` also relies on external tool definitions which may have their own configurations.

**Edge Cases/Launch Conditions:**

The functions will fatally log and exit if:
1.  Configuration loading fails.
2.  Tool execution fails.
3.  The external tool returns an empty result.

There are no explicit command-line arguments or flags defined in this snippet, but the configuration itself may be loaded via such mechanisms elsewhere. The package assumes that `tools.ToolsExecutor` is properly initialized and configured before use.

**Project Package Structure:**

```
examples/codemonkey/pkg/planner/
├── planner.go
```

**Code Relations & Unclear Places:**

The core logic revolves around using external tools (`tools.ReWOOToolDefinition`) to process string inputs into executable actions. The exact behavior depends entirely on the configuration of these tools and the prompts provided in `lePromptGithelper` and `lePromptCLI`.  There's no explicit error handling beyond fatal logging, which suggests a lack of recovery mechanisms.

The source of the input strings for `PlanGitHelper` and `PlanCLIExecutor` is not defined within this snippet; they are likely passed from another part of the system. The purpose of these functions (e.g., what kind of "Reviewer Result" or "task" they process) remains unclear without further context.

<end_of_output>
```