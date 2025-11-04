# pkg/config/calloptions.go  
## Package / Component    
**Name:** `config`    
  
### Imports    
- `github.com/tmc/langchaingo/llms`  
  
---  
  
## External Data & Input Sources    
The function `ConifgToCallOptions` consumes a single argument of type `DefaultCallOptions`.    
All fields that may be present in this struct are examined and, if non‑nil, converted into an `llms.CallOption`:  
  
| Field | Purpose | llms option used |  
|-------|---------|-------------------|  
| `Model` | Model to use for the call | `WithModel` |  
| `CandidateCount` | Number of candidates to generate | `WithCandidateCount` |  
| `MaxTokens` | Max tokens per candidate | `WithMaxTokens` |  
| `Temperature` | Temperature setting (note: typo in name) | `WithTemperature` |  
| `StopWords` | Stop‑word list | `WithStopWords` |  
| `TopK` | Top‑k sampling parameter | `WithTopK` |  
| `TopP` | Top‑p sampling parameter | `WithTopP` |  
| `Seed` | Random seed for reproducibility | `WithSeed` |  
| `MinLength` | Minimum length of generated text | `WithMinLength` |  
| `MaxLength` | Maximum length of generated text | `WithMaxLength` |  
| `N` | Number of completions to request | `WithN` |  
| `RepetitionPenalty` | Penalty for repeated tokens | `WithRepetitionPenalty` |  
| `FrequencyPenalty` | Frequency penalty setting | `WithFrequencyPenalty` |  
| `PresencePenalty` | Presence‑penalty setting (typo in name) | `WithPresencePenalty` |  
| `JSONMode` | Boolean flag to enable JSON mode | `WithJSONMode` |  
| `ResponseMIMEType` | MIME type for the response | `WithResponseMIMEType` |  
  
---  
  
## TODOs    
No explicit `TODO:` comments are present in this file.  
  
---  
  
## Summary of Major Code Parts    
  
### 1. Function `ConifgToCallOptions`  
*Purpose:* Convert a `DefaultCallOptions` struct into a slice of `llms.CallOption`s that can be passed to an LLM call.    
*Logic Flow:*  
- Initialise an empty slice `opts`.  
- For each field in the input struct, check if it is non‑nil.  
- Append the corresponding `llms.With…` option to `opts`.    
  - The function covers all available fields, ensuring that only set values are added.  
- Return the populated slice.  
  
*Key Points:*  
- Handles optional configuration elegantly by checking for nil pointers before appending.  
- Uses a consistent pattern (`if cfg.Field != nil { opts = append(opts, llms.With…(*cfg.Field)) }`) which makes it easy to extend or modify.  
- The final `return opts` provides the complete set of options ready for use elsewhere in the package.  
  
---  
  
# pkg/config/config.go  
## Package / Component    
**Name:** `config`    
  
### Imports    
```go  
import (  
	"fmt"  
	"os"  
	"reflect"  
	"regexp"  
	"strconv"  
	"strings"  
  
	"github.com/joho/godotenv"  
	"github.com/rs/zerolog/log"  
)  
```  
  
### External Data / Input Sources    
* Environment variables read from a `.env` file (via `godotenv.Load()`) and the OS environment (`os.LookupEnv`, `os.Environ`).    
* The code expects an optional prefix key defined by the constant `EnvPrefixKey`.    
  
### TODOs    
No explicit `TODO:` comments are present in this file.    
  
---  
  
## Summary of Major Code Parts  
  
| Section | Purpose & Key Details |  
|---------|-----------------------|  
| **`Config` struct** | Holds top‑level configuration values for AI URL, token, model, and a nested `DefaultCallOptions`. It also contains several feature flags (`ReWOODisable`, `SemanticSearchDisable`, etc.) and maps/slices that use wildcard tags (e.g., `CommandExecutorCommands`). |  
| **`DefaultCallOptions` struct** | Nested configuration for AI call options. Each field has an `env:"..."` tag that may contain multiple keys separated by commas, allowing fallback to alternative env names. |  
| **`NewConfig()`** | Creates a new `Config`, loads the `.env` file, applies any prefix from `EnvPrefixKey`, and iterates over all struct fields calling `processField`. It validates required top‑level values (`AIURL`, `AIToken`, `Model`). |  
| **`processField(...)`** | Handles a single field: builds the full env key (prefix + tag), determines if it is a nested struct, map or slice, and delegates to `processWildcardField` when the tag ends with `_*`. It then sets the field value via `setFieldValue`. |  
| **`processWildcardField(...)`** | Special logic for fields whose tag contains a wildcard (`_ *`). It supports maps and slices: it builds a regex to match env keys like `<prefix>_<suffix>` and populates either a map or slice accordingly. |  
| **`setFieldValue(...)`** | Generic setter that interprets the kind of the field (string, int, bool, float64/32, slice) and assigns the parsed value from the environment variable. It also handles pointer dereferencing when needed. |  
  
---  
  
The file implements a flexible configuration loader that can read nested structs, maps, and slices from environment variables with optional prefixes and wildcard keys. All logic is contained in this single package, making it straightforward to extend or modify.  
  
