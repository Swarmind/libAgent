# main

This package is a command-line tool that handles user input, executes tools, and processes results. It is part of the examples directory and includes multiple files.

## Package Files
- `main.go`: Contains the main entry point for the example.
- `examples/hacker/main.go`: Another entry point for a different example.

## External Data/Config Options
- **Config File**: Loaded via `config.NewConfig()`; must be valid.
- **Config Fields**:
  - `SemanticSearchDisable` (bool): Disables semantic search in the tool.

## Behavior
The package's behavior is affected by the configuration, such as disabling semantic search.

## Edge Cases
- If user input is empty, the tool may fail (no handling for this).
- If `toolsExecutor.CallTool` returns an error, the program exits with a fatal log.
- If the tool returns an empty result, the program exits with a fatal log.

## Assumptions
- The tool `tools.ReWOOToolDefinition.Name` is available and properly configured.
- The tool's response is expected to be a non-empty string.

## Notes
- No TODOs or comments in the code.