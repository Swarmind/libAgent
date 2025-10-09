## Package Summary: `main`

This package serves as the entry point for a command-line application that leverages external planners and executors to process natural language input into executable commands. It's part of an example project (`github.com/Swarmind/libagent/examples/codemonkey`). The code currently focuses on CLI execution with hardcoded input, but contains commented-out sections suggesting original intent for GitHub issue processing via webhooks.

**Imports:**

*   `fmt`: Standard Go package for formatted I/O.
*   `github.com/Swarmind/libagent/examples/codemonkey/pkg/executor`: Custom package responsible for executing commands (likely shell or system calls).
*   `github.com/Swarmind/libagent/examples/codemonkey/pkg/planner`: Custom package that translates natural language input into executable plans.

**Configuration:**

The commented-out webhook handler suggests the use of environment variables like `LISTEN_ADDR` for configuring the server address, but this functionality is currently disabled. The main function uses a hardcoded prompt string: `"Find current OS version and OS type (windows/linux/android)"`. No other configuration options are present in the active code.

**Execution Flow:**

1.  The `main` function calls `planner.PlanCLIExecutor` with the hardcoded prompt to generate an executable plan.
2.  It then prints the generated plan and executes it using `executor.ExecuteCommands`.

**Inactive Functionality (Commented-Out):**

*   A webhook handler (`es.StartWebhookHandler`) intended for receiving GitHub issue events is disabled. This section suggests processing issues, removing tags, and passing the result to a planner.
*   A test snippet demonstrates manual invocation of `reviewer.GatherInfo`, planning with `planner.PlanGitHelper`, and printing results; this functionality is also inactive.

**Project Package Structure:**

```
examples/codemonkey/cmd/codemonkey/
├── main.go
```

The package consists solely of the `main.go` file, which contains the entry point for the application. The external packages (`pkg/executor`, `pkg/planner`) are assumed to be defined elsewhere within the project.