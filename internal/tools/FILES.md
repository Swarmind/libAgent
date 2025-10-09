# internal/tools/tools.go  
## Tools Package Summary  
  
**Package Name:** `tools`  
  
**Imports:**  
  
*   `context`: For managing request contexts.  
*   `encoding/json`: For JSON serialization and deserialization.  
*   `fmt`: For formatted I/O.  
*   `slices`: For slice manipulation (sorting).  
*   `strings`: For string operations.  
*   `github.com/rs/zerolog/log`: For structured logging.  
*   `github.com/tmc/langchaingo/llms`: LangChain LLM definitions and interfaces.  
  
**External Data / Inputs:**  
  
*   The package relies on external `ToolData` structures, which contain function definitions (`llms.FunctionDefinition`), call functions (`func(context.Context, string) (string, error)`), and optional cleanup functions (`func() error`).  
*   It accepts a context (`context.Context`) for tool execution.  
*   The package takes `llms.ToolCall` structs as input to execute tools with arguments.  
  
**TODOs:**  
  
No TODO comments found in the provided code snippet.  
  
### Code Sections Summary:  
  
**1. LLM Definition:**  
  
Defines a default LLM function definition (`LLMDefinition`) for general-purpose language model interactions, including its name and description. This is used as a base tool within the executor.  
  
**2. ToolData Struct:**  
  
Represents a single tool with its definition, execution logic (Call), and optional cleanup function.  This struct encapsulates all necessary components to interact with an external functionality.  
  
**3. ToolsExecutor Struct & Methods:**  
  
The core component for managing and executing tools. It stores a map of `ToolData` indexed by tool name (`Tools`). Key methods include:  
  
*   `Execute`: Executes a given tool call, retrieves the result, and returns it as a `llms.ToolCallResponse`.  
*   `GetTool`: Retrieves a specific tool from the internal map.  
*   `CallTool`: Calls the execution function of a specified tool with provided arguments.  
*   `ToolsList`: Returns a sorted list of available tools in the format expected by LangChain's LLM interface.  
*   `ToolsPromptDesc`: Generates a human-readable description of all available tools, including their names, input types, and descriptions for use in prompts.  This is crucial for providing context to an LLM about its capabilities.  
*   `ProcessToolCalls`: Executes multiple tool calls sequentially, handling errors and collecting the results into a single string.  
*   `Cleanup`: Iterates through all tools and executes their cleanup functions if defined.  
  
**4. Utility Functions:**  
  
The code includes utility functions for sorting tools by name (`slices.SortFunc`) and formatting tool descriptions using JSON marshaling to extract input types from function definitions.  These ensure consistent presentation of available tools.  
  
