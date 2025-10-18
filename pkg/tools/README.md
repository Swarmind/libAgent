```markdown
# libagent/pkg/tools

This package provides a collection of tools for use within the libagent framework, including search (DuckDuckGo), command execution, Metasploit integration, Nmap scanning, semantic search via PostgreSQL and OpenAI, ReWOO reasoning engine, and web content extraction.  Tools are registered globally during initialization based on configuration settings.

## Project Package Structure:

```
pkg/tools/
├── ddgSearch.go
├── executor.go
├── metasploit.go
├── nmap.go
├── rewoo.go
├── semanticSearch.go
├── tools.go
└── webReader.go
```

### Configuration and Environment Variables:

*   **`config.Config`**:  The primary configuration source for all tools, loaded from external sources (e.g., files, environment variables). Key settings include disabling individual tools (`DDGSearchMaxResults`, `MsfDisable`, `NmapDisable`, `ReWOODisable`, `SemanticSearchDBConnection`, `WebReaderDisable`).
*   **`PATH`**: Used by the command executor for shell execution.

### Command-Line Arguments/Edge Cases:

Most tools rely on JSON input via an LLM interface, but some have configuration dependencies that can prevent them from loading if not properly set up (e.g., missing database credentials for semantic search). The command executor is particularly sensitive to environment variables and requires a functional shell (`bash`) to operate correctly.  The Metasploit tool depends on the `msfconsole` executable being available in the system's PATH.

### Relations Between Code Entities:

*   **`tools.go`**: Acts as the central registry for all tools, managing their initialization and enabling/disabling based on configuration.
*   **Individual Tool Files (`ddgSearch.go`, `executor.go`, etc.)**: Each file defines a specific tool with its own logic for processing input, executing commands (if applicable), and returning results.  They all adhere to the LangChain-compatible function definition pattern.
*   **`internal/tools`**: Contains shared internal implementations used by multiple tools (e.g., `webReader`).

### Unclear Places or Dead Code:

The comment in `nmap.go` about moving `GenerateMsfQueries` suggests potential refactoring but doesn't indicate immediate issues. The "should NOT exist arguments which called NAME cause it cause COLLISION with actual function name" comment in `semanticSearch.go` is vague and requires further investigation to determine if it represents a real problem or just poor naming convention.  The global singleton pattern for the tool executor (`globalToolsExecutor`) might introduce testability challenges but isn't inherently problematic.
```