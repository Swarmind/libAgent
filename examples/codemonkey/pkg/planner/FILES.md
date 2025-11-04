# examples/codemonkey/pkg/planner/planner.go  
**Package name:** `planner`  
  
### Imports  
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
- `context` – standard Go context handling  
- `encoding/json` – JSON marshaling for tool queries  
- `os` – OS interaction (stderr output)  
- `config` – local configuration loader  
- `tools` – ReWOOTool executor utilities  
- `zerolog` & `log` – structured logging  
  
### External data / input sources  
| Source | Purpose |  
|--------|---------|  
| `lePromptGithelper` | Prompt template for Git helper tool |  
| `lePromptCLI` | Prompt template for CLI executor tool |  
| `config.NewConfig()` | Loads configuration needed by the tools executor |  
| `tools.ReWOOToolDefinition.Name` | Tool name to whitelist in executor |  
| `tools.NewToolsExecutor(...)` | Creates an executor instance |  
| `json.Marshal(rewooQuery)` | Serializes query payload for ReWOO tool |  
  
### TODOs  
- No explicit TODO comments, but the two plan functions are identical except for prompt source; consider extracting common logic into a helper.  
  
---  
  
## Summary of major code parts  
  
#### 1. Prompt constants  
```go  
const lePromptGithelper = `Role: You are an Instruction Synthesis Agent...`  
const lePromptCLI = `You are an AI command generation assistant specialized in creating executable CLI command sequences...`  
```  
These strings contain the instruction templates used to build queries for the ReWOO tool.  
  
#### 2. `PlanGitHelper(review string) string`  
- Sets global log level and logger.  
- Loads configuration via `config.NewConfig()`.  
- Creates a context (`context.Background()`).  
- Whitelists the ReWOOToolDefinition.Name in an executor.  
- Builds a query by concatenating `lePromptGithelper` with the provided review text.  
- Marshals the query to JSON, calls the tool, and returns its result.  
  
#### 3. `PlanCLIExecutor(task string) string`  
- Mirrors `PlanGitHelper`, but uses `lePromptCLI` as the prompt source.  
- All other steps (config load, executor creation, query build, call, return) are identical.  
  
Both functions can be refactored into a single generic helper that accepts a prompt prefix; this would reduce duplication and improve maintainability.  
  
