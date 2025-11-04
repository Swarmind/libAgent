# internal/tools/rewoo/rewoo.go  
**Package / Component name**    
`rewoo`  
  
---  
  
## Imports  
| Package | Purpose |  
|---------|---------|  
| `context` | Go context handling for async calls |  
| `encoding/json` | JSON marshaling/unmarshaling of tool results |  
| `fmt` | String formatting |  
| `regexp` | Regular expression parsing of plan steps |  
| `slices` | Slice utilities (used to check if a key exists) |  
| `strings` | String manipulation |  
| `github.com/Swarmind/libagent/internal/tools` | Tool executor and helper functions |  
| `github.com/Swarmind/libagent/pkg/util` | Utility helpers (`RemoveThinkTag`) |  
| `github.com/google/uuid` | UUID generation for decision markers |  
| `graph "github.com/JackBekket/langgraphgo/graph/stategraph"` | State‑graph builder & runner |  
| `github.com/rs/zerolog/log` | Structured logging |  
| `github.com/tmc/langchaingo/llms` | LLM interface (chat, calls) |  
| `github.com/tmc/langchaingo/llms/openai` | OpenAI LLM implementation |  
  
---  
  
## External data / input sources  
* **LLM** – an instance of `openai.LLM` is used to generate content for each step.  
* **ToolsExecutor** – provides tool descriptions, calls and cleanup logic.  
* **StateGraph** – the workflow graph that connects plan → tool → solve nodes.  
* **Prompts** – four large prompt constants (`PromptGetPlan`, `PromptSolver`, `PromptDecision`, `PromptRegeneratePlan`) drive the LLM interactions.  
  
---  
  
## TODO comments  
No explicit `TODO:` markers are present in this file.    
(If future work is needed, a dedicated section can be added.)  
  
---  
  
# Summary of major code parts  
  
### 1. Prompt constants  
Four multiline string constants define the language templates used by the LLM:  
* **PromptGetPlan** – instructs the LLM to produce step‑by‑step plans with evidence variables.  
* **PromptSolver** – feeds the generated plan into a solver prompt and expects a final answer.  
* **PromptDecision** – evaluates whether the solved plan is correct, using a UUID marker for verification.  
* **PromptRegeneratePlan** – regenerates a new plan if the previous one was unsatisfactory.  
  
### 2. Type definitions  
| Type | Description |  
|------|-------------|  
| `ReWOO` | Main orchestrator struct holding an LLM instance, a tools executor and default call options. |  
| `State` | Runtime state of the workflow: current attempt, task description, plan string, steps, results map, solved plan text, final answer. |  
| `Step` | One step of the plan – contains plan text, name, tool to use, and tool input. |  
  
### 3. Regular expression  
* **StepPattern** – captures each plan line:    
  `Plan:<text> #E<name>=<tool>[<input>]`.    
  It is used in `GetPlan` and `ObserveEnd` to split the LLM output into individual steps.  
  
### 4. Graph initialization (`InitializeGraph`)  
Builds a state‑graph with three nodes:  
* `"plan"` – calls `GetPlan`  
* `"tool"` – calls `ToolExecution`  
* `"solve"` – calls `Solve`  
  
Edges: plan → tool, solve → observeEnd (conditional), tool → route (conditional). Entry point is the plan node.  
  
### 5. Plan extraction (`GetPlan`)  
1. If no plan string exists, it asks the LLM to generate one using `PromptGetPlan` and the list of available tools.  
2. Parses all matches with `StepPattern`, stores them in a map keyed by evidence variable (e.g., `#E1`).  
3. Orders steps according to the keys and appends them into `state.Steps`.  
4. Logs the extracted steps.  
  
### 6. Solve step (`Solve`)  
Iterates over each step, replaces placeholders with previously obtained results, builds a consolidated plan string, then asks the LLM to solve the task using `PromptSolver`.    
The final answer is stored in `state.Result`.  
  
### 7. Tool execution (`ToolExecution`)  
* Picks the current step based on how many results have already been produced.  
* Builds a prompt for the chosen tool (LLM or other) with `PromptCallTool`.  
* If the tool is not an LLM, it fetches its description via `ToolsExecutor.GetTool` and optionally prepares a command‑executor query.  
* Executes the tool call through the LLM, processes returned tool calls, marshals JSON into `state.Results`.  
  
### 8. Routing (`Route`)  
Decides whether to go back to the solve node or to the tool node depending on how many results have been produced.  
  
### 9. Observe end (`ObserveEnd`)  
* Checks if maximum attempts were reached; otherwise it asks the LLM to decide if the plan is correct using `PromptDecision`.  
* If the decision marker appears in the response, the workflow ends.  
* Otherwise a new plan is generated with `PromptRegeneratePlan`, state is reset for another attempt, and the graph loops back to the plan node.  
  
### 10. Helper (`getCurrentTask`)  
Returns the index of the next step to execute based on how many results are already present.  
  
---  
  
**Overall flow**    
`rewoo` orchestrates a multi‑step LLM workflow: it generates a plan, executes each tool call in order, solves the task, and verifies the result. The graph ensures that after each solve or tool execution the next appropriate node is visited until the maximum number of attempts (`ObserveAttempts`) is reached.  
  
