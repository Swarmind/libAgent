<end_of_output>

# `rewoo`

`internal/tools/rewoo/rewoo.go` implements a lightweight orchestration layer for an LLM‑driven workflow that

1. **generates a step‑by‑step plan** (`PromptGetPlan`);
2. **executes each tool call** in the order defined by the plan;
3. **solves the task** with the accumulated plan (`PromptSolver`);
4. **verifies and optionally regenerates** the plan until a maximum number of attempts is reached.

The package exposes a single public struct `ReWOO` that holds an LLM instance, a tools executor, and default call options.  
All logic is driven by a small state‑graph with three nodes (`plan`, `tool`, `solve`) and two conditional edges (`solve → observeEnd`, `tool → route`).  

---

## Project structure

```
internal/tools/rewoo/
├── rewoo.go
```

---

## Environment variables, flags & command‑line arguments

| Variable / Flag | Purpose |
|------------------|---------|
| `ObserveAttempts` (int) | Maximum number of plan iterations before the workflow stops. |
| `TaskDescription` (string) | Human‑readable description that is fed to the LLM in all prompts. |
| `LLM` (`openai.LLM`) | The OpenAI instance used for all chat calls. |
| `ToolsExecutor` (`tools.Executor`) | Helper that fetches tool descriptions and executes them. |

These values are read from the struct fields of `ReWOO`.  
If the package is built as a CLI, they can be supplied via flags such as:

```
--attempts <int>
--task "<string>"
--llm-config "<json>"
```

---

## Edge cases for launching

* **CLI entry point** – If this package contains a `main()` function (or is imported by one), it can be started with `go run internal/tools/rewoo/main.go` after setting the above variables.
* **Library usage** – Other packages may instantiate `ReWOO`, set its fields, and call `InitializeGraph()` to kick off the workflow.  
  The graph will automatically loop back to the plan node if the final answer is not yet satisfactory.

---

## Summary of major code parts

1. **Prompt constants** – four multiline strings (`PromptGetPlan`, `PromptSolver`, `PromptDecision`, `PromptRegeneratePlan`) that form the LLM templates for each step.
2. **Type definitions** –  
   * `ReWOO` – orchestrator struct.  
   * `State` – runtime state (attempt, task, plan string, steps, results map, solved text).  
   * `Step` – a single plan line with name, tool and input.
3. **Regular expression** – `StepPattern` extracts each plan line (`Plan:<text> #E<name>=<tool>[<input>]`) into a map keyed by the evidence variable.
4. **Graph initialization** – `InitializeGraph()` builds nodes `"plan"`, `"tool"` and `"solve"` with edges that form the execution loop.
5. **GetPlan()** – generates or parses a plan, orders steps, logs them, and stores them in `state.Steps`.
6. **Solve()** – iterates over all steps, replaces placeholders, builds a consolidated plan string, asks the LLM to solve it, and records the final answer.
7. **ToolExecution()** – selects the current step based on how many results exist, builds a prompt for that tool, executes it via the LLM, and marshals JSON into `state.Results`.
8. **Route()** – decides whether to go back to the solve node or to the tool node depending on progress.
9. **ObserveEnd()** – checks if the maximum number of attempts has been reached; otherwise it verifies the plan with `PromptDecision` and may regenerate a new plan with `PromptRegeneratePlan`.

---

## Relations between code entities

* The graph nodes call the methods in the order: `plan → GetPlan`, `tool → ToolExecution`, `solve → Solve`.  
  After each solve, `observeEnd()` decides whether to finish or loop back.  
* `State` is passed through all callbacks; it holds the current attempt counter (`Attempt`) and a map of results (`Results`).  
* The regular expression in `GetPlan()` feeds directly into the ordering logic used by `Solve()`.  
* `Route()` uses the length of `state.Results` to determine whether another tool call or another solve is needed.

---

## Possible dead code

No unused functions or variables were detected; all exported methods are referenced by the graph edges.