```markdown
## rewoo Package Summary

This package implements a reasoning and tool execution workflow using LangChainGo, designed to solve tasks by generating plans, executing tools based on those plans, and refining the process if necessary. It leverages LLMs (specifically OpenAI) for planning, solving, and decision-making. The core logic revolves around a state graph that orchestrates these steps.

**Imports:**

*   `context`: For managing request contexts.
*   `encoding/json`: For handling JSON data in tool interactions.
*   `fmt`: For formatted string output.
*   `regexp`: For parsing plans using regular expressions.
*   `slices`: For slice operations (e.g., checking if an element exists).
*   `strings`: For string manipulation.
*   `github.com/Swarmind/libagent/internal/tools`: Custom tools executor for external tool interactions.
*   `github.com/Swarmind/libagent/pkg/util`: Utility functions, including tag removal from LLM responses.
*   `github.com/google/uuid`: For generating unique identifiers (e.g., decision markers).
*   `graph "github.com/JackBekket/langgraphgo/graph/stategraph"`: State graph implementation for workflow orchestration.
*   `github.com/rs/zerolog/log`: Structured logging.
*   `github.com/tmc/langchaingo/llms`: LangChainGo LLM interface.
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI LLM implementation.

**External Data & Inputs:**

*   The package relies on an external `ToolsExecutor` to manage and execute tools (e.g., search, calculator). The tool definitions are dynamically loaded from this executor.
*   Input tasks are provided as strings (`State.Task`).
*   LLM responses are used throughout the workflow for planning, solving, and decision-making.

**TODOs:**

No explicit `TODO` comments were found in the code.

### Core Components:

#### 1.  Workflow Graph Initialization (`InitializeGraph`)

The package sets up a state graph with nodes for "plan," "tool," and "solve." Edges define the execution flow, including conditional transitions based on decision-making (e.g., regenerating plans if necessary). The entry point is the "plan" node.

#### 2.  Plan Generation (`GetPlan`)

The `GetPlan` function uses an LLM to generate a step-by-step plan for solving the given task, including tool calls with input parameters. It parses the generated plan using regular expressions to extract individual steps (tool name, input). The extracted steps are stored in the `State.Steps`.

#### 3.  Tool Execution (`ToolExecution`)

The `ToolExecution` function executes tools based on the current step in the plan. It prepares tool calls with appropriate arguments and invokes the external `ToolsExecutor`. Tool outputs are captured and stored in the `State.Results` map for use in subsequent steps. The LLM is used to sanitize arguments before execution.

#### 4.  Solving (`Solve`)

The `Solve` function uses an LLM to solve the task based on the generated plan and executed tool results. It combines the plan, evidence (tool outputs), and task description into a prompt for the LLM. The final answer is stored in `State.Result`.

#### 5.  Decision Making & Plan Regeneration (`ObserveEnd`)

The `ObserveEnd` function evaluates whether the solved plan is correct using an LLM-based decision process. If incorrect, it regenerates the plan and restarts the workflow. This loop continues until a satisfactory solution is found or a maximum number of attempts is reached. The regeneration uses a prompt to fix errors in the previous plan.

#### 6.  Routing (`Route`)
The `Route` function determines whether to proceed with tool execution or solving based on the current task state. If all tools have been executed, it transitions to the "solve" node; otherwise, it continues executing tools.

**Constants:**

*   `GraphPlanName`, `GraphToolName`, `GraphSolveName`: Node names in the workflow graph.
*   `ObserveAttempts`: Maximum number of plan regeneration attempts.
*   `PromptGetPlan`, `PromptSolver`, `PromptDecision`, `PromptRegeneratePlan`, `PromptLLMTool`, `PromptCallTool`: LLM prompts used for different stages of the workflow.
```