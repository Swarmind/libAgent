# Package **executor**

## Project structure
```
examples/codemonkey/pkg/executor/
└─ executor.go
```

---

## Short overview  
`executor.go` implements a small command‑line helper that builds an AI prompt, calls a tool executor and runs the resulting shell script.  
The package pulls in standard Go libraries for context handling, JSON encoding, OS interaction and string manipulation, plus two local packages from *libagent* (`config`, `tools`) and the third‑party logging library **zerolog**.

---

## Environment variables / flags / command‑line arguments  

| Variable / flag | Purpose | Default value |
|------------------|---------|---------------|
| `TASK` (env)    | Task description passed to the generator. | *none* – must be supplied by caller |
| `-task string`  | CLI flag for the same task. | *required* |

The package can be invoked directly as a Go binary:

```bash
go run examples/codemonkey/pkg/executor/executor.go -task "build docker image"
```

or built into an executable and used in a larger workflow.

---

## Key functions & their interactions  

| Function | Responsibility |
|----------|----------------|
| `CliGenerator(task string) string` | 1. Sets log level to **Debug**.<br>2. Creates a new configuration instance (`cfg`).<br>3. Builds a whitelist of tool names that will be used by the executor.<br>4. Instantiates a `ToolsExecutor` with the provided context, config and whitelist.<br>5. Constructs a query for the *ReWOOTool* using `CreatePrompt(task)`, marshals it to JSON, and calls the tool via the executor.<br>6. Returns the raw string result from the tool call. |
| `CreatePrompt(task string) string` | Builds a human‑readable prompt that instructs an AI assistant to generate Unix/Linux CLI commands for the supplied task. The prompt contains rules, examples, and the actual task, then returns it as a single string ready for JSON marshalling. |
| `ExecuteCommands(commandStr string) error` | Splits the raw command string into individual lines, trims whitespace, filters out empty entries.<br>Joins the cleaned commands back together with newline separators to form a shell script.<br>Executes the script using `sh -c` via the OS exec package.<br>Prints the combined output of the script to standard error for debugging. |

The three functions are tightly coupled:  
* `CliGenerator` relies on `CreatePrompt`, and `ExecuteCommands` can be used to run the generated commands.

---

## How the application can be launched  

1. **Direct CLI** – as shown above, passing a task string via the `-task` flag.  
2. **As part of a larger workflow** – import the package in another Go file and call `executor.CliGenerator("…")`.  
3. **Scheduled job** – wrap the binary in a cron or CI pipeline; environment variable `TASK` can be set instead of CLI flag.

---

## Summary of logic  

1. The program starts by creating a configuration instance (`cfg`) that holds runtime settings for the executor.  
2. It builds a whitelist of tool names (currently two: *ReWOOTool* and *CommandExecutorDefinition*) that will be used to generate commands.  
3. A prompt is created with `CreatePrompt`, which contains instructions for an AI model to produce shell commands for the supplied task.  
4. The prompt is marshalled into JSON, sent to the tool executor, and the resulting command string is returned by `CliGenerator`.  
5. Finally, `ExecuteCommands` can be called to split that string into a script, execute it with `sh -c`, and log the output.

The package therefore acts as a thin wrapper around an AI‑powered command generator, turning a task description into executable shell commands.

---