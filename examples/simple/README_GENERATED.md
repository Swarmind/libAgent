# simple

## Overview
A simple example demonstrating the usage of the `libagent` library to integrate with OpenAI LLMs. It provides a basic workflow for logging, configuration loading, LLM initialization, agent execution, and result handling.

## Structure
- `main.go`: Contains the entry point and core logic for the example.

## Key Functionality
1. **Logging**: Sets up debug-level logging to stderr with console formatting.
2. **Configuration**: Loads settings via `config.NewConfig()` (see below for details).
3. **LLM Integration**: Initializes an OpenAI LLM client using configured values.
4. **Agent Execution**: Runs a simple agent with a hardcoded test prompt.
5. **Result Handling**: Outputs processed results after removing control tags.

## Configuration Options
The following environment variables/fields are required via `config.NewConfig()`:
- `AIURL` (string): OpenAI API base URL
- `AIToken` (string): API authentication token
- `Model` (string): Model identifier to use
- `DefaultCallOptions` (config.CallOptions): Default call parameters for LLM

## Issues & Limitations
- **Typo**: `config.ConifgToCallOptions` (should be `Config`)
- **No validation**: No checks for empty/invalid AIURL/AIToken/Model values
- **No nil handling**: Potential nil responses from `agent.SimpleRun()` are not addressed
- **Static Prompt**: Uses a hardcoded test string ("This is a test...") instead of dynamic input
- **Limited Demonstration**: Only shows basic LLM interaction, not advanced agent capabilities

## Usage
Run with appropriate environment variables set for the configuration options.