# main

This package is a benchmarking tool that tests different models by executing a set of tasks. It iterates over a list of 20 model identifiers, runs a multi-step task plan, and handles errors and edge cases.

## Package File Structure
- `main.go` (in `examples/` directory): Contains the main logic for benchmarking.
- `main.go` (in `examples/modelsbenchmark/` directory): Part of the same `main` package, likely containing specialized benchmarking functionality.

## Description
The package is designed to evaluate model behavior by executing a predefined set of tasks (e.g., file operations, port checks, Git operations) against a list of 20 model identifiers. It uses a `ReWOO` tool to execute tasks and handles errors by logging and sleeping for 2 minutes (except for the last model).

## Configuration/External Data
- **Configuration**: `config.NewConfig()` (loads external configuration, not shown in code).
- **Tools**: Supports `ReWOOTool` and `CommandExecutor` (whitelisted).
- **ModelList**: 20 model identifiers (e.g., "qwen3-32b", "gemma-3-27b-it").
- **Prompt**: 11 sequential tasks (file operations, port checks, Git operations, etc).

## Edge Cases/Notes
- The **last model** in `ModelList` does **not** trigger a 2-minute sleep after errors.
- `ToolsExecutor` is **recreated and cleaned** for each model.
- `ReWOO` is **executed for every model**.
- The code **exits on fatal errors** (e.g., config loading failure).
- **Non-OK status** from `ReWOO` is **not handled**.
- **LocalAI watchdog** is handled via 2-minute sleeps (explicitly coded).