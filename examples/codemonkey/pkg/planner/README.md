# Package `planner`

## Overview
`planner.go` implements two helper functions that build and execute ReWOO queries for Git‑helper and CLI‑executor tools.  
Both helpers share identical logic – only the prompt template differs – so a single generic routine could be extracted.

---

## File structure

```
examples/codemonkey/pkg/planner/
└── planner.go
```

* `planner.go` – main source file for the package.

---

## Imports & dependencies
```go
import (
	"context"
	"encoding/json"

	"os"

	"github.com/Swarmind/libagent/pkg/config"
	"github.com/Swarmind/libagent/pkg/tools"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)
```
| Import | Purpose |
|--------|---------|
| `context` | Go context handling for tool execution |
| `encoding/json` | Marshal query payloads to JSON |
| `os` | Write logs / errors to stderr |
| `config` | Load configuration needed by the tools executor |
| `tools` | ReWOOTool executor utilities |
| `zerolog` & `log` | Structured logging |

---

## Environment variables, flags and CLI arguments
* **Environment** – none explicitly read; the package relies on defaults from `config.NewConfig()`.
* **Flags / command‑line args** – not present in this file; the two exported functions accept a single string argument each (`review` or `task`).  
  They can be invoked directly from another package (e.g. a CLI entry point) as:
  ```go
  result := planner.PlanGitHelper("my review")
  ```
* **Files & paths** – only `planner.go`; no sub‑packages.

---

## External data / input sources

| Source | Purpose |
|--------|---------|
| `lePromptGithelper` | Prompt template for the Git helper tool |
| `lePromptCLI` | Prompt template for the CLI executor tool |
| `config.NewConfig()` | Loads configuration needed by the tools executor |
| `tools.ReWOOToolDefinition.Name` | Tool name to whitelist in executor |
| `json.Marshal(rewooQuery)` | Serializes query payload for ReWOO tool |

---

## Core logic

### Prompt constants
```go
const lePromptGithelper = `Role: You are an Instruction Synthesis Agent...`
const lePromptCLI = `You are an AI command generation assistant specialized in creating executable CLI command sequences...`
```

These strings form the instruction prefix for each ReWOO query.

### `PlanGitHelper(review string) string`
1. Sets global log level and logger (`zerolog.SetGlobalLevel(zerolog.InfoLevel)`).
2. Loads configuration via `config.NewConfig()`.
3. Creates a background context.
4. Whitelists the tool name in an executor instance.
5. Builds a query by concatenating the Git‑helper prompt with the supplied review text.
6. Marshals the query to JSON, calls the tool, and returns its result.

### `PlanCLIExecutor(task string) string`
Identical flow as above but uses the CLI prompt (`lePromptCLI`) instead of the Git helper one.

---

## Relations & potential refactor
* Both functions perform the same steps; only the prompt source differs.  
  A single generic routine such as:
  ```go
  func planWithPrompt(prompt string, reviewOrTask string) string { … }
  ```
  could replace both and reduce duplication.

* The package currently has no dead code – all imports are used, and every function is exported for external use.

---

## Edge cases & launch scenarios
* **Library usage** – other packages can import `planner` and call either helper directly.
* **CLI entry point** – a separate main package could expose a command that forwards user input to one of these helpers and prints the result.  
  Example:
  ```bash
  go run ./cmd/planner-cli --mode=githelper --review="…"
  ```
  (the actual CLI flags would be defined in the main package, not shown here).

---

## Summary
`planner.go` provides two convenience functions that generate ReWOO queries for Git‑helper and CLI‑executor tools.  
The code is straightforward; the only improvement would be to merge the duplicated logic into a single helper.