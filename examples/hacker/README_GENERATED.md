# hacker

## Overview
This package provides a command-line interface for interacting with the ReWOO tool, handling user input, and executing tasks through a configured tools executor. It is part of a larger system involving configuration management and tool execution.

## File Structure
- `main.go`: Contains the entry point and core logic for user interaction, tool invocation, and result handling.

## Behavior Control
### External Configuration Options
- `cfg.SemanticSearchDisable` (bool): Disables semantic search in the configuration. Default: `false`.

## Key Details
- **ReWOO Tool**: Invoked with a fixed `Prompt` constant; its purpose and `ReWOOToolArgs` structure are not documented.
- **Configuration**: Initializes a config object with `SemanticSearch` disabled by default.
- **Tool Execution**: Uses `tools.NewToolsExecutor` and handles cleanup via `defer`.

## Issues & Notes
- No validation for empty user input (scanner.Scan() could fail).
- No handling for invalid JSON marshaling (errors are logged but not addressed).
- `toolsExecutor.Cleanup()` is deferred but errors are not checked beyond logging.
- The `Prompt` constant is hard-coded and cannot be modified via configuration.
- The initial comment mentions "manual use of specific tool" but does not explain the ReWOO tool's purpose or `ReWOOToolArgs` expectations.