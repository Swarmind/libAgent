# tools

## Overview
The `tools` package provides a framework for managing and executing tools, with a focus on integrating with Large Language Models (LLMs). It enables the registration, execution, and cleanup of tools, including a special global `LLMDefinition` for general-purpose AI model interactions.

## File Structure
- **tools.go**: Core implementation of the `ToolsExecutor` and related types (e.g., `ToolData`, `LLMDefinition`).
- **rewoo/rewoo.go**: Unspecified functionality (part of the package but not detailed here).
- **webReader/webReader.go**: Unspecified functionality (part of the package but not detailed here).

## Key Components
- **LLMDefinition**: A global `llms.FunctionDefinition` representing a default LLM. Cannot be modified.
- **ToolData**: Encapsulates a tool's definition, execution logic, and cleanup handler.
- **ToolsExecutor**: Manages a registry of tools and provides methods to execute, call, list, and process tool interactions.

## External Configuration
- **LLMDefinition**: A global constant; cannot be modified.
- **Tools Registration**: Tools are registered via the `Tools` map in `ToolsExecutor`.
- **ToolData Requirements**:
  - `FunctionDefinition` (Name, Description, Parameters).
  - `Call` handler: `func(context.Context, string) (string, error)`.
  - Optional `Cleanup` handler: `func() error`.

## Notes on Edge Cases & Issues
- **LLMDefinition**: Not part of the `Tools` map but treated as a special case in `ToolsPromptDesc`.
- **ToolsPromptDesc**: Complex parameter handling; may panic on invalid input.
- **ProcessToolCalls**: Only returns the last tool's response (overwrites prior results).
- **Sorting Logic**: Duplicated in `ToolsList` and `ToolsPromptDesc`.
- **Safety**: `ToolsExecutor` does not validate if a `Call` function is `nil` (may panic).
- **Cleanup**: Stops on the first error, but spec is unclear on whether all cleanup should proceed.