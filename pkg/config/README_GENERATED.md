# pkg/config

## Overview
The `pkg/config` package handles configuration management for application settings and LLM (Large Language Model) call behavior. It provides mechanisms to load environment variables into structured configurations and convert configuration options into LLM call parameters.

## File Structure
- **calloptions.go**: Converts `DefaultCallOptions` to `llms.CallOption` slices for LLM interactions.
- **config.go**: Loads and validates environment variables into a `Config` struct, handling complex configuration scenarios.

## External Configuration Options
The package uses environment variables to configure behavior. Key options include:

### Core Application Configuration
- `AI_URL`, `AI_TOKEN`, `MODEL`: Required fields for API access and model selection.
- `REWO_DISABLE`, `REWO_DEFAULT_CALL_OPTIONS`: Controls Rewoo integration behavior.
- `SEMANTIC_SEARCH_DISABLE`, `DDG_SEARCH_DISABLE`: Enables/disables semantic search and DuckDuckGo integration.
- `WEB_READER_DISABLE`, `NMAP_DISABLE`: Disables web reading or network scanning features.
- `COMMAND_EXECUTOR_DISABLE`, `COMMAND_EXECUTOR_COMMANDS`: Manages command execution via environment variables.

### LLM Call Behavior (via `DefaultCallOptions`)
The following fields are mapped to environment variables (e.g., `TEMPERATURE` → `Temperature`):
- `Model`, `CandidateCount`, `MaxTokens`, `Temperature`, `StopWords`, `TopK`, `TopP`, `Seed`, `MinLength`, `MaxLength`, `N`, `RepetitionPenalty`, `FrequencyPenalty`, `PresencePenalty`, `JSONMode`, `ResponseMIMEType`  
*(18 optional fields; only non-nil values are applied)*

## Notes & Issues
- **Typo**: Function `ConifgToCallOptions` should be `ConfigToCallOptions`.
- **No Documentation**: No comments for functions/parameters in `calloptions.go`.
- **No Validation**: Invalid values (e.g., non-integer strings for `int` fields) are not handled.
- **Limited Wildcard Support**: Only `map`/`slice` types are supported for wildcard fields (e.g., `COMMAND_EXECUTOR_CMD_*`).
- **Unresolved Dependencies**: `DefaultCallOptions` is referenced but not defined in this package.
- **No Error Handling**: Critical fields (e.g., `AIURL`) are validated, but invalid values in non-critical fields are ignored.