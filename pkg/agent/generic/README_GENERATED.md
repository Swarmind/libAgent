# generic

## Overview
The `generic` package provides a foundational agent implementation that integrates Large Language Models (LLMs) with tool execution capabilities. It manages interactions between LLMs (currently OpenAI) and tools via a `ToolsExecutor`, enabling structured handling of chat-based workflows and tool invocations.

## Key Features
- **LLM Integration**: Uses OpenAI LLMs for generating responses via `LLM.GenerateContent`.
- **Tool Management**: Dynamically loads and executes tools through `ToolsExecutor` (requires `ToolsList()` and `ProcessToolCalls()` implementations).
- **Dual Execution Modes**: 
  - `Run()` handles chat histories as `MessageContent` sequences.
  - `SimpleRun()` simplifies interactions by accepting raw strings.

## Configuration & Behavior
### External Controls
- **LLM Behavior**: Configured via `llms.CallOption` (e.g., temperature, max tokens).
- **Tool Handling**: Requires `ToolsExecutor` to implement:
  - `ToolsList()`: Returns available tools (loaded once lazily).
  - `ProcessToolCalls()`: Handles tool execution results.

### Limitations
- **No Safety Checks**: Assumes `response.Choices[0]` is non-empty/valid.
- **Static Tools**: `toolsList` is not reinitialized if `ToolsExecutor.ToolsList()` changes after initialization.
- **Error Handling Gaps**: No explicit handling for `ToolsExecutor.ToolsList()` or `ProcessToolCalls()` failures.
- **Code Duplication**: `Run()` and `SimpleRun()` share significant logic.

## Issues
- **Unaddressed Edge Cases**: 
  - Undefined behavior if `LLM.GenerateContent` returns no valid response.
  - Tools cannot be modified after initial load.
- **Refactoring Opportunities**: 
  - Shared logic between `Run()` and `SimpleRun` could be consolidated.
  - Additional safety checks for LLM responses are missing.