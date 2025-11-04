# Package config

The **config** package provides a lightweight, environment‑driven configuration loader for an LLM‑based application.  
It reads values from a `.env` file (or any other source that `godotenv.Load()` can read) and turns them into a strongly typed struct that can be passed to the rest of the codebase.

---

## File structure

```
pkg/config/
├── calloptions.go
└── config.go
```

* **calloptions.go** – helper that converts a `DefaultCallOptions` value into a slice of `llms.CallOption`s.
* **config.go** – main configuration struct, loader logic and helpers for handling nested structs, maps and slices.

---

## Core data structures

| File | Type | Purpose |
|------|------|---------|
| `calloptions.go` | `func ConifgToCallOptions(cfg *DefaultCallOptions) []llms.CallOption` | Builds a list of LLM call options from the nested struct. |
| `config.go` | `type Config struct { … }` | Holds top‑level config values (`AIURL`, `AIToken`, `Model`) and a nested `DefaultCallOptions`. |
|  | `type DefaultCallOptions struct { … }` | Individual LLM call parameters (model, candidate count, temperature, etc.). |

---

## Environment variables & flags

The loader expects the following keys in the environment (or `.env` file).  
All names are optional; if a key is missing it will simply be ignored.

| Key | Source | Description |
|-----|--------|-------------|
| `AIURL` | top‑level | URL of the LLM endpoint. |
| `AIToken` | top‑level | API token for authentication. |
| `Model` | top‑level | Name of the model to use. |
| `DefaultCallOptions_*` | nested | Each field in `DefaultCallOptions` can be supplied via a key that ends with `_ *`. The loader will match any suffix after the underscore, e.g. `DefaultCallOptions_Model`, `DefaultCallOptions_CandidateCount`, etc. |

The struct tags (`env:"..."`) allow multiple names per field; the loader picks the first one that exists.

---

## How it works

1. **`NewConfig()`** – called by external code to create a fully populated `Config`.  
   * Loads `.env` via `godotenv.Load()`.  
   * Iterates over all fields of `Config`, calling `processField` for each.  
   * Validates that required top‑level values (`AIURL`, `AIToken`, `Model`) are present.

2. **`processField()`** – generic handler for any field in the struct.  
   * Builds a full key name by concatenating an optional prefix (`EnvPrefixKey`) with the tag value.  
   * Detects whether the field is a nested struct, map or slice and delegates to `processWildcardField` when the tag ends with `_ *`.  
   * Calls `setFieldValue()` to actually assign the parsed value.

3. **`processWildcardField()`** – special logic for fields that contain a wildcard (`_*`).  
   * Uses a regular expression to match all environment keys that start with the prefix and end with any suffix.  
   * Supports both maps (e.g., `CommandExecutorCommands`) and slices, filling them accordingly.

4. **`setFieldValue()`** – low‑level helper that interprets the kind of the field (`string`, `int`, `bool`, `float64/32`, slice) and assigns the parsed value from the environment variable.  
   * Handles pointer dereferencing when needed (e.g., for nested structs).  

5. **`ConifgToCallOptions()`** – turns a populated `DefaultCallOptions` into a slice of `llms.CallOption`s that can be passed to an LLM call.  
   * Checks each field; if it is non‑nil, appends the corresponding `With…` option from the `github.com/tmc/langchaingo/llms` package.

---

## Edge cases & launch scenarios

* **Missing keys** – If a key is absent, the loader simply skips that field.  
* **Multiple names per field** – The tag can contain several comma‑separated names; the first one found will be used.  
* **Running the application** – A typical usage pattern would be:  

  ```bash
  go run ./cmd/main.go          # main package imports pkg/config
  ```

  or, if this is a library only, simply import `pkg/config` in any other Go file and call `config.NewConfig()` to obtain a ready‑to‑use configuration.

---

## Summary

The **config** package offers:

* A clean, extensible way to read environment variables into a typed struct.  
* Automatic handling of nested structs, maps and slices via wildcard tags.  
* A helper that converts the nested `DefaultCallOptions` into LLM call options ready for use in the rest of the application.

This makes it straightforward to change configuration values without touching code – just update the `.env` file or set environment variables before launching the program.