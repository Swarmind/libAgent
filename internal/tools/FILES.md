# internal/tools/tools.go  
# Package & Imports  
**Package name:** `tools`    
**Imports:**  
```go  
import (  
	"context"  
	"encoding/json"  
	"fmt"  
	"slices"  
	"strings"  
  
	"github.com/rs/zerolog/log"  
	"github.com/tmc/langchaingo/llms"  
)  
```  
  
---  
  
## External Data & Input Sources  
- **LLMDefinition** – a global `llms.FunctionDefinition` that describes the base “LLM” tool.    
  It contains a name, description and will be used as part of the prompt description.  
  
---  
  
## TODOs  
No explicit `TODO:` comments are present in this file.  
  
---  
  
## Struct Definitions  
| Name | Purpose |  
|------|---------|  
| `ToolData` | Holds a single tool’s definition, its execution function (`Call`) and an optional cleanup routine. |  
| `ToolsExecutor` | Manages a map of tools (`map[string]*ToolData`). Provides methods to execute calls, list tools, build prompt descriptions, process multiple calls, and clean up resources. |  
  
---  
  
## Methods Overview  
  
### `Execute`  
Executes a single tool call:  
1. Builds a response header from the incoming `llms.ToolCall`.    
2. Calls `CallTool` to obtain the raw content string.    
3. Populates the response body and returns it.  
  
### `GetTool`  
Retrieves a `*ToolData` by name from the internal map, returning an error if not found.  
  
### `CallTool`  
Convenience wrapper that fetches the tool data via `GetTool` and invokes its `Call` function with the supplied arguments.  
  
### `ToolsList`  
Creates a slice of `llms.Tool` objects (type “function”) for all registered tools, sorted alphabetically by name.    
Useful when building a list of available tools for an LLM prompt.  
  
### `ToolsPromptDesc`  
Builds a human‑readable description string that lists every tool’s name, input type and description.    
It starts with the global `LLMDefinition`, appends all registered tools, sorts them, then formats each entry as:  
```  
(idx) ToolName[inputType]: Description  
```  
  
### `ProcessToolCalls`  
Iterates over a slice of incoming calls, executes each via `Execute`, logs any errors and collects the final content string.    
The current implementation overwrites `content` on every iteration; it could be extended to concatenate results.  
  
### `Cleanup`  
Runs the optional cleanup function for each tool in the map, propagating any error that occurs.  
  
---  
  
## Summary  
This file defines a lightweight executor (`ToolsExecutor`) that manages a collection of LLM‑based tools.    
It exposes helper methods to list available tools, generate prompt descriptions, execute individual calls and process batches of calls.    
The global `LLMDefinition` provides a base tool definition that is automatically included in the prompt description.  
  
