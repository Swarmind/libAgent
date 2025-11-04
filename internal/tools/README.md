# Package `tools`

## Quick Summary  
`internal/tools` implements a lightweight LLM‑tool executor that can be used by any higher‑level package to register, list and invoke “function” tools for a language‑model prompt. The core type is **`ToolsExecutor`**, which keeps a map of `ToolData` objects (each containing a name, description, call function and optional cleanup).  Methods on the executor build prompt descriptions, execute single or multiple calls, and expose the tool list to an LLM.

---

## File structure

```
internal/tools/
├─ rewoo/rewoo.go          # (not shown – likely contains helper for “rewoo” tool)
├─ tools.go                 # core executor logic
└─ webReader/webReader.go  # (not shown – probably a helper to read from the web)
```

---

## Environment variables / flags / cmd‑line arguments  
No explicit environment variables or command‑line flags are defined in this package.  All configuration is done through the exported methods of `ToolsExecutor`.

---

## How the code works

| Entity | Purpose |
|--------|---------|
| **`ToolData`** | Holds a single tool’s definition (`llms.FunctionDefinition`), its execution function (`Call`) and an optional cleanup routine. |
| **`ToolsExecutor`** | Keeps a map `map[string]*ToolData`.  It offers:  
  * `Execute(ctx, call)` – runs one tool call.  
  * `GetTool(name)` – fetches a tool by name.  
  * `CallTool(ctx, name, args…)` – convenience wrapper around `GetTool`.  
  * `ToolsList()` – returns an alphabetically sorted slice of all tools as `llms.Tool` objects.  
  * `ToolsPromptDesc()` – builds a human‑readable description string for the LLM prompt (global definition + all registered tools).  
  * `ProcessToolCalls(ctx, calls…)` – iterates over a slice of incoming calls and executes each via `Execute`.  
  * `Cleanup()` – runs cleanup functions for every tool. |

The executor is initialized with a global `llms.FunctionDefinition` called **`LLMDefinition`** that represents the base “LLM” tool; this definition is automatically included in the prompt description.

---

## Edge cases & launch scenarios  

* **CLI/command‑line usage** – If another package imports `internal/tools`, it can create a new `ToolsExecutor`, register tools, and call `Execute()` or `ProcessToolCalls()` from a CLI command.  
* **Main package integration** – A main program could instantiate the executor, populate its map (e.g., via `tools.go`), then use `ToolsList()` to feed an LLM prompt and finally invoke `ProcessToolCalls()` after receiving tool calls from the model.  

---

## Summary of logic flow

1. **Initialization** – Create a `ToolsExecutor`, optionally pre‑populate it with tools (e.g., via `rewoo` or `webReader`).  
2. **Prompt building** – Call `ToolsPromptDesc()` to get a string that lists all tool names, input types and descriptions; feed this into an LLM prompt.  
3. **Execution** – After the model returns calls, use `ProcessToolCalls()` to run each call through the executor’s map.  
4. **Cleanup** – When finished, invoke `Cleanup()` to release any resources.

The package is intentionally minimal: all state lives in the executor map, and the only external dependency is `github.com/tmc/langchaingo/llms`.  The code can be extended by adding more tools or improving error handling, but as‑is it already supports a full tool lifecycle.