# examples/codemonkey/pkg/executor/executor.go  
# Package / Component    
**executor**  
  
## Imports  
```go  
import (  
    "context"  
    "encoding/json"  
    "fmt"  
    "os"  
    "os/exec"  
    "strings"  
  
    "github.com/Swarmind/libagent/pkg/config"  
    "github.com/Swarmind/libagent/pkg/tools"  
  
    "github.com/rs/zerolog"  
    "github.com/rs/zerolog/log"  
)  
```  
The package pulls in standard Go libraries for context handling, JSON encoding, OS interaction and string manipulation.    
External dependencies are the local `config` and `tools` packages from *libagent*, plus the third‑party logging library **zerolog**.  
  
## External Data / Input Sources  
| Source | Purpose |  
|--------|---------|  
| `config.NewConfig()` | Loads configuration data for the executor. |  
| `tools.WithToolsWhitelist(...)` | Supplies a list of tool names that will be used by the executor. |  
| `tools.ReWOOToolDefinition.Name` & `tools.CommandExecutorDefinition.Name` | Identifiers for the two tools that are whitelisted and later invoked. |  
  
## TODOs  
No explicit TODO comments were found in this file.  
  
---  
  
## Summary of Major Code Parts  
  
### 1. `CliGenerator(task string) string`  
* Sets global log level to **Debug** and configures a console logger.  
* Creates a new configuration instance (`cfg`) and a background context.  
* Builds a whitelist of tool names that will be used by the executor.  
* Instantiates a `ToolsExecutor` with the provided context, config, and whitelist.  
* Constructs a query for the *ReWOOTool* using `CreatePrompt(task)`, marshals it to JSON, and calls the tool via the executor.  
* Returns the raw string result from the tool call.  
  
### 2. `CreatePrompt(task string) string`  
* Builds a human‑readable prompt that instructs an AI assistant to generate Unix/Linux CLI commands for the supplied task.  
* The prompt contains rules, examples, and the actual task, then returns it as a single string ready for JSON marshalling.  
  
### 3. `ExecuteCommands(commandStr string) error`  
* Splits the raw command string into individual lines, trims whitespace, and filters out empty entries.  
* Joins the cleaned commands back together with newline separators to form a shell script.  
* Executes the script using `sh -c` via the OS exec package.  
* Prints the combined output of the script to standard error for debugging.  
  
---  
  
All functions are tightly coupled: `CliGenerator` relies on `CreatePrompt`, and `ExecuteCommands` can be used to run the generated commands. The code is ready for integration into a larger *libagent* workflow.  
  
