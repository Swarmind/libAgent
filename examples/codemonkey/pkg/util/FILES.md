# examples/codemonkey/pkg/util/utils.go  
**Package / Component Name**    
`utility`  
  
---  
  
### Imports  
| Package | Purpose |  
|---------|---------|  
| `os` | Provides access to environment variables via `os.Getenv`. |  
| `github.com/rs/zerolog/log` | Enables structured logging; used for fatal error reporting. |  
  
---  
  
### External Data / Input Sources  
* The function expects a key name (`key string`) that identifies an environment variable.  
* It reads the value of this variable using `os.Getenv`.  
* If the retrieved value is empty, it logs a fatal message indicating the missing key.  
  
---  
  
### TODOs  
No explicit TODO comments are present in the file.  
  
---  
  
## Summary of Major Code Parts  
  
### Function `GetEnv`  
- **Signature**: `func GetEnv(key string) string`  
- **Behavior**:  
  - Calls `os.Getenv` to fetch the value associated with the supplied key.  
  - Checks if the returned value is empty; if so, logs a fatal message using `log.Fatal().Msgf`.  
  - Returns the fetched value for further use by callers.  
  
This single function encapsulates environment variable retrieval and basic error handling, making it a convenient helper within the `utility` package.  
  
