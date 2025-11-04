# tools

The **`tools`** package bundles a collection of LLM‑exposed utilities that can be invoked from higher‑level workflows or directly via the command line.  
Each file implements one tool, registers it in a global registry, and the package also provides a lightweight executor that stitches all registered tools together.

---

## File structure

```
pkg/tools/
├─ attacker.go
├─ ddgSearch.go
├─ executor.go
├─ metasploit.go
├─ nmap.go
├─ rewoo.go
├─ semanticSearch.go
├─ tools.go
├─ webReader.go
└─ windows_executor.go
```

---

## Environment variables / configuration flags

| Flag (in `config.Config`) | Purpose |
|------------------------------|---------|
| `ExploitDisable` | Enable/disable the Metasploit attacker tool. |
| `DDGSearchDisable` | Enable/disable DuckDuckGo search wrapper. |
| `CommandExecutorDisable` | Enable/disable the generic command executor (both Linux and Windows). |
| `ReWOODisable` | Enable/disable the rewoo reasoning graph tool. |
| `SemanticSearchDisable` | Enable/disable semantic vector‑store search. |
| `WebReaderDisable` | Enable/disable web‑reader. |
| `NmapDisable` | Enable/disable nmap scanner. |

Additional flags used by individual tools:

* `DDGSearchMaxResults`, `DDGSearchUserAgent`
* `SemanticSearchAIURL`, `SemanticSearchAIToken`, `SemanticSearchDBConnection`,
  `SemanticSearchEmbeddingModel`, `SemanticSearchMaxResults`
* `CommandExecutorCommands` (map of command → description)

---

## How the package can be launched

1. **CLI** – a main binary that imports `pkg/tools` will call `tools.NewToolsExecutor(ctx, cfg, opts...)`.  
2. The executor iterates over all init functions in `globalToolsRegistry`, builds a map of tool data keyed by function name, and exposes them via the LLM interface (`llms.FunctionDefinition`).  
3. From the CLI or another component you can invoke any registered tool by its name (e.g., `"exploit"`, `"webSearch"`, `"commandExecutor"`).  

---

## Tool relationships & key logic

| File | Main role | Key entities |
|------|------------|--------------|
| `attacker.go` | Metasploit module runner (`exploit`) | `ExploitToolDefinition`, `ExploitToolArgs`, `(*ExploitTool).Call` |
| `ddgSearch.go` | DuckDuckGo search wrapper (`webSearch`) | `DDGSearchDefinition`, `DDGSearchArgs`, `DDGSearchTool` |
| `executor.go` | Generic command executor (`commandExecutor`) | `CommandExecutorDefinition`, `CommandExecutorArgs`, `(*CommandExecutorTool).Call` |
| `metasploit.go` | Metasploit query generator & runner (`msf_search`) | `MsfSearchToolDefinition`, `MsfSearchToolArgs`, `GenerateMsfQueries` |
| `nmap.go` | Nmap scanner + Metasploit query builder (`nmap`) | `NmapToolDefinition`, `NmapToolArgs`, `ParseNmapPorts` |
| `rewoo.go` | OpenAI‑based reasoning graph executor (`rewoo`) | `ReWOOToolDefinition`, `ReWOOToolArgs`, `(*ReWOOTool).Call` |
| `semanticSearch.go` | Vector‑store semantic search (`semanticSearch`) | `SemanticSearchDefinition`, `SemanticSearchArgs`, `SemanticSearchTool` |
| `tools.go` | Executor factory & registry helper | `NewToolsExecutor`, `WithToolsWhitelist` |
| `webReader.go` | URL → markdown reader (`webReader`) | `WebReaderDefinition`, `WebReaderArgs`, `(*WebReaderTool).Call` |
| `windows_executor.go` | Windows CMD executor (`windowsCommandExecutor`) | `WCommandExecutorDefinition`, `WCommandExecutorArgs`, `(*WCommandExecutorTool).Call` |

All tools register themselves via an `init()` that appends a closure to `globalToolsRegistry`.  
The executor collects these closures, applies optional whitelist filtering, and builds a map of `tools.ToolData` objects.  

---

## Edge cases & launch options

* **Whitelist** – Pass `WithToolsWhitelist("exploit", "webSearch")` when creating the executor to limit which tools are available.
* **Default values** – Each tool checks its config flag; if disabled it simply returns `nil`.  
  * For example, `ddgSearch.go` sets defaults for max results and user agent if zero/empty.  
* **Command lists** – The generic executors (`executor.go`, `windows_executor.go`) read a map of commands from the config and append them to their description; this allows the LLM to see available sub‑commands.
* **Error handling** – All init functions return an error; if any fails, the executor will skip that tool but continue building others.

---

## Summary

The `tools` package provides:

1. A set of small, self‑contained tools (Metasploit, DuckDuckGo, Nmap, OpenAI reasoning, web reader, etc.).
2. Each tool declares an LLM function definition and a call method that parses JSON input, performs its work, and returns raw output.
3. An executor factory (`tools.go`) that aggregates all registered tools into a single `ToolsExecutor` instance.
4. Configuration is driven by the global `config.Config`; flags control enable/disable status and provide defaults for each tool.

This package can be used as a library in a larger agent or run directly via a CLI that creates an executor, then calls any of the exposed functions from the LLM interface.