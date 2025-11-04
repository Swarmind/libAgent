# Package / Component **reviewer**

The `reviewer` package implements a small helper that gathers information from a GitHub issue and returns it as a string. It is part of the *reviewer* package, which orchestrates calls to external tools (ReWOOTool, SemanticSearch, DDGSearch) via a `ToolsExecutor`.

## File structure
```
examples/codemonkey/pkg/reviewer/
└── reviewer.go
```

## Environment variables / configuration flags

| Variable | Purpose | Default |
|----------|---------|---------|
| `LOG_LEVEL` | Sets the logging level for the executor (used in `reviewer.go`) | `debug` |
| `CONFIG_PATH` | Path to a JSON/YAML config file that contains runtime settings for the executor | `config.yaml` |
| `TOOL_LIST` | Comma‑separated list of tool names to be used by the executor | `"ReWOOToolDefinition,DDGSearchDefinition,SemanticSearchDefinition"` |

## Command‑line arguments

The package can be invoked as a command line tool or imported into another Go program. If it is built as a CLI:

```
go run ./examples/codemonkey/pkg/reviewer --issue="123" --repo="my-repo"
```

* `--issue` – GitHub issue number to gather information for.
* `--repo` – Repository name that will be used in the prompt.

## Edge cases / launch scenarios

| Scenario | How to launch |
|----------|---------------|
| **CLI** | `go run ./examples/codemonkey/pkg/reviewer --issue=123 --repo=my-repo` |
| **Binary** | `go build -o reviewer ./examples/codemonkey/pkg/reviewer && ./reviewer --issue=123 --repo=my-repo` |
| **Library import** | Import the package in another module: `import "github.com/Swarmind/libagent/pkg/reviewer"` and call its exported functions. |

## Summary of major code parts

### `GatherInfo`
* Sets global log level to debug and configures the logger.
* Loads configuration, creates a context, and builds a whitelist of tools.
* Instantiates a `ToolsExecutor` with the given context, config, and tool list.
* Builds a query struct (`ReWOOToolArgs`) using `CreatePrompt`, marshals it into JSON, and calls the ReWOO tool via the executor.
* Returns the raw string result from the tool call.

### `CreatePrompt`
* Constructs a detailed prompt for the ReWOO tool.  
  * The prompt instructs the AI to research an issue in a repository, using semantic search and DDG search tools.  
  * It includes placeholders for the issue title, repo name, and the names of the two tools that will be used.  
* Returns the formatted string ready for JSON serialization.

The file is self‑contained: it pulls together configuration, tool execution, and prompt generation to produce a single string result that can be consumed by other components in the *reviewer* package.