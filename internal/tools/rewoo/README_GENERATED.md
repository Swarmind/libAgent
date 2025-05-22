# rewoo

This package provides a workflow management system for task execution using LLMs and tools. It handles task planning, execution, and solving through a state graph.

**External Configuration:**
- Tools registered with ToolsExecutor (e.g., search, LLM, Calculator)
- LLM API endpoint (OpenAI) via `openai.LLM`
- `DefaultCallOptions` (LLM invocation settings)
- Prompt templates (e.g., `PromptGetPlan`, `PromptSolver`)

**Notes:**
- The `StepPattern` regex may not handle all edge cases (e.g., malformed plans)
- The code assumes the last match in the plan is the valid one (due to map overwriting)
- No error handling for invalid tool inputs in `ToolExecution`
- The `GraphPlanName/GraphToolName/GraphSolveName` constants are used for graph node management

**Edge Cases:**
- Empty plan generation (returns error)
- Mismatch between steps and results (e.g., incomplete execution)
- Invalid tool inputs (e.g., non-existent tools in the plan)
- Large/complex plans that exceed LLM output limits