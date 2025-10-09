## Config Package Summary

This package manages application configuration through environment variables, `.env` files, and reflection-based field population. It provides a structured way to define settings for AI integrations, search functionality, command execution, and LLM call options. The core logic revolves around the `Config` struct, which is populated from environment variables using tags (`env`) that specify variable names. Wildcard matching (`_*`) allows dynamic loading of map or slice values based on prefixed environment variables.

**Configuration Sources:**

*   `.env` files (optional): Loaded if present to provide default configuration values.
*   Environment Variables: Primary source, accessed via `os.Getenv()`. Supports prefixing with `LIBAGENT_ENV_PREFIX`. Wildcard support for map/slice types using the `_*` suffix in env tags.

**Key Structures:**

*   `Config`: Main configuration struct holding AI settings (URL, token, model), search options, tool disabling flags, and command execution configurations (`CommandExecutorCommands`).
*   `DefaultCallOptions`: Nested structure defining default LLM call parameters (model, max tokens, temperature, etc.).

**Functions:**

*   `NewConfig()`: Initializes a `Config` instance by loading environment variables and populating fields using reflection. Validates essential AI settings (URL, token, model).
*   `processField()`: Recursive function that iterates through struct fields, extracts the `env` tag, retrieves the corresponding environment variable value, and sets it to the field. Handles nested structs recursively.
*   `processWildcardField()`: Processes wildcard-tagged fields (`_*`) by extracting all matching environment variables and populating either a map or slice based on the field type.
*   `setFieldValue()`: Sets the value of a struct field using reflection, handling string, int, bool, float64/32, and slice types with necessary conversions.

**Environment Variables:**

The package relies heavily on environment variables prefixed by `LIBAGENT_ENV_PREFIX` (if defined). Wildcard matching (`_*`) allows loading multiple values into maps or slices based on the prefix. All fields within both `Config` and `DefaultCallOptions` structs are populated from these variables if present. The absence of required AI settings (URL, token, model) will not cause an error but may lead to undefined behavior in dependent components.

**Project Package Structure:**

```
pkg/config/
├── calloptions.go
└── config.go
```

**Relations between Code Entities:**

The `Config` struct contains a nested `DefaultCallOptions` struct, which is populated independently from environment variables using the same reflection-based mechanism. The `NewConfig()` function orchestrates the entire process by recursively calling `processField()` to populate all fields in both structs. Wildcard processing (`processWildcardField()`) extends this functionality to handle dynamic map/slice configurations.