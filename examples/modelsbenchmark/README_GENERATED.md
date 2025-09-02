# modelsbenchmark

## Overview
A benchmarking tool for evaluating model behavior using the ReWOO framework. Executes tests against a hardcoded list of models, focusing on model response handling and resource management.

## File Structure
- `main.go`: Contains the entry point and core benchmarking logic.

## Package Functionality
- Benchmarks models by executing the ReWOO tool with a fixed prompt
- Supports model testing through the `ModelList` (hardcoded) and `Prompt` (static)
- Manages tool execution via `toolsExecutor` with ReWOO and CommandExecutor tools

## Configuration & Behavior
**External Data/Config Options:**
- `config.NewConfig()`: Loads configuration (not shown in code)
- `cfg.Model`: Dynamically set to each model in `ModelList`
- `toolsExecutor`: Uses `ReWOOTool` and `CommandExecutor` (whitelisted)

**Behavior Notes:**
- Model testing is non-interactive (fixed prompt and model list)
- Includes a 2-minute sleep between tests (unexplained in code)
- `toolsExecutor.Cleanup()` is called after each iteration (purpose unclear)

## Issues & Concerns
- **Invalid Models**: Contains unsupported model names (e.g., "mlabonne_qwen3-8b-abliterated")
- **Unclear Logic**: 
  - `toolsExecutor.Cleanup()` purpose is undefined
  - `ReWOOTool` response handling is undocumented
  - `time.Sleep(2 * time.Minute)` is used for LocalAI watchdog behavior (not explained)
- **Inefficiency**: `toolsExecutor` is recreated for each model (could be optimized)

## Status
This is a prototype benchmarking implementation with several unresolved edge cases and potential optimizations.