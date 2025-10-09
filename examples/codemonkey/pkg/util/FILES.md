# examples/codemonkey/pkg/util/utils.go  
## Package: `utility`  
  
This package provides utility functions, specifically for retrieving environment variables with error handling.  
  
**Imports:**  
  
*   `os`: For accessing environment variables.  
*   `github.com/rs/zerolog/log`: For logging fatal errors if an environment variable is missing.  
  
**External Data / Input Sources:**  
  
The primary input to this package is the operating system's environment variables, accessed via `os.Getenv()`.  It relies on these being properly set before execution.  
  
**Function Summary:**  
  
### `GetEnv(key string) string`  
  
This function retrieves the value of an environment variable specified by the `key`. If the environment variable is not set (empty string), it logs a fatal error message using `zerolog` and terminates the program. Otherwise, it returns the retrieved value as a string. This ensures that critical configuration values are present before proceeding.  
  
**TODOs:**  
  
There are no TODO comments in this code snippet.  
  
