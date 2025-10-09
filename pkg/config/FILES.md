# pkg/config/calloptions.go  
## Package: `config`  
  
This package provides functionality for converting configuration settings into LangChainGo call options. It primarily focuses on mapping fields from a `DefaultCallOptions` struct to the appropriate `llms.CallOption` parameters used by the LangChainGo library.  
  
**Imports:**  
  
*   `github.com/tmc/langchaingo/llms`:  Used for accessing and defining LLM call options.  
  
**External Data / Input Sources:**  
  
The function `ConifgToCallOptions` takes a single input:  
  
*   `cfg DefaultCallOptions`: A struct containing optional configuration parameters for the LLM call (model, candidate count, max tokens, temperature, stop words, topK, topP, seed, minLength, maxLength, N, repetition penalty, frequency penalty, presence penalty, JSON mode, and response MIME type).  All fields within this struct are pointers to allow for optional configuration.  
  
**Function Summary:**  
  
*   `ConifgToCallOptions(cfg DefaultCallOptions) []llms.CallOption`: This function converts the provided `DefaultCallOptions` into a slice of `llms.CallOption`. It iterates through each configurable field in the input struct and, if present (not nil), appends the corresponding LangChainGo call option to the resulting slice. The function returns an empty slice if no configuration parameters are set.  
  
**TODOs:**  
  
There are no TODO comments within this code file.  
  
# pkg/config/config.go  
## Config Package Summary  
  
**Package Name:** `config`  
  
**Imports:**  
  
*   `fmt`: For formatted I/O.  
*   `os`: For interacting with the operating system (environment variables).  
*   `reflect`: For runtime reflection.  
*   `regexp`: For regular expressions.  
*   `strconv`: For string conversions.  
*   `strings`: For string manipulation.  
*   `github.com/joho/godotenv`: For loading environment variables from `.env` files.  
*   `github.com/rs/zerolog/log`: For structured logging.  
  
**External Data Sources:**  
  
*   `.env` file (optional): Loaded if present, providing default values for configuration parameters.  
*   Environment Variables: Primary source of configuration data, accessed via `os.Getenv()`.  Supports prefixing environment variables using the `LIBAGENT_ENV_PREFIX` variable. Wildcard support is available for map and slice types with the `_*` suffix in env tags.  
  
**TODOs:** None found within this file.  
  
### Config Struct  
  
The core configuration structure, `Config`, holds various settings related to AI integrations (URL, token, model), search functionality (semantic, DDG), and tool disabling flags (ReWOO, web reader, Nmap, MSF).  It also includes a map for command execution configurations (`CommandExecutorCommands`). The struct uses environment variables with the `env` tag to populate its fields.  
  
### DefaultCallOptions Struct  
  
A nested structure used within `Config` to define default options for AI calls (model, candidate count, max tokens, temperature, stop words, top-k/p sampling parameters, seed, length constraints, repetition penalties, JSON mode, response MIME type).  Fields are also populated from environment variables.  
  
### NewConfig Function  
  
This function initializes a `Config` instance by loading environment variables and populating the struct fields using reflection. It handles optional `.env` file loading, prefixing of environment variable names based on `LIBAGENT_ENV_PREFIX`, wildcard matching for map/slice types, and performs basic validation (ensuring AI URL, token, and model are set).  
  
### processField Function  
  
This recursive function iterates through the fields of a struct using reflection. It extracts the `env` tag from each field, constructs the environment variable name based on the prefix (if any), retrieves the value from the environment, and sets it to the corresponding field in the config struct. Handles nested structs recursively. Supports wildcard matching for map/slice types with the `_*` suffix in env tags.  
  
### processWildcardField Function  
  
Handles processing of fields tagged with a wildcard (`_*`). It extracts all environment variables matching the prefix defined by the tag, and populates either a map or slice based on the field type. If the key is an integer it will be used as index for slice population.  
  
### setFieldValue Function  
  
Sets the value of a struct field using reflection. Supports string, int, bool, float64/32, and slice types. Performs necessary conversions (e.g., `strconv.Atoi`, `strconv.ParseBool`) and handles pointer dereferencing if needed. Returns an error if parsing fails or the type is unsupported.  
  
