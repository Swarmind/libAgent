# generic

## Overview
`pkg/agent/generic/agent.go` implements a lightweight OpenAI‑based agent that can execute LLM calls and process tool invocations.  
The package exposes two public methods:

| Method | Purpose |
|--------|---------|
| `Run` | Execute an LLM request with a pre‑built state slice, automatically adding the configured tools and returning the first AI response. |
| `SimpleRun` | Same as `Run`, but accepts a single string input instead of a message slice – useful for quick one‑shot calls. |

Both methods rely on the `openai.LLM` implementation from *tmc/langchaingo* and a `tools.ToolsExecutor` that handles tool calls returned by the model.

---

## File structure

```
pkg/
└── agent/
    └── generic/
        └── agent.go
```

Only one source file is present; all logic lives in `agent.go`.

---

## Key code entities and their relationships

| Entity | Description |
|--------|-------------|
| `Agent` struct | Holds the LLM instance, a tools executor, and a pointer to the list of tools (`toolsList`). |
| `Run(ctx context.Context, state []llms.MessageContent, opts ...llms.CallOption)` | 1. Ensures `toolsList` is non‑nil. <br>2. Adds the tool list to the LLM options via `llms.WithTools`. <br>3. Calls `a.LLM.GenerateContent(ctx, state, opts...)`. <br>4. Processes any returned tool calls with `a.ToolsExecutor.ProcessToolCalls`. <br>5. Returns the first choice’s content as a chat message. |
| `SimpleRun(ctx context.Context, input string, opts ...llms.CallOption)` | Same flow as `Run`, but wraps the single string into a human‑type chat message before passing it to the LLM. |

The two methods share most of their logic; only the way the initial state is built differs.

---

## Configuration knobs

| Variable / flag | Where it appears | Usage |
|------------------|-------------------|-------|
| `a.LLM` | `Agent` struct | Must be set to an initialized `*openai.LLM`. |
| `a.ToolsExecutor` | `Agent` struct | Handles tool calls; can be configured externally. |
| `a.toolsList` | `Agent` struct | List of tools passed to the LLM via `llms.WithTools`. |
| `opts ...llms.CallOption` | Both methods | Additional options such as temperature, max tokens, etc., that are forwarded to `GenerateContent`. |

No explicit environment variables or command‑line flags exist in this file; configuration is done programmatically by setting the struct fields before calling either method.

---

## How to launch / use

1. **As a library** – import `"github.com/Swarmind/libagent/pkg/agent/generic"` and create an `Agent` instance:

   ```go
   ag := &generic.Agent{
       LLM:           openai.NewLLM(...),
       ToolsExecutor: tools.NewToolsExecutor(),
       toolsList:     &[]llms.Tool{...},
   }
   ```

2. **CLI / main package** – create a `main.go` that imports this package and calls either method:

   ```go
   func main() {
       ctx := context.Background()
       ag := generic.NewAgent(...) // helper constructor can be added later
       resp, err := ag.SimpleRun(ctx, "Hello world", llms.WithTemperature(0.7))
       fmt.Println(resp)
   }
   ```

Edge cases:
- If `toolsList` is nil when calling `Run`, the method will panic; ensure it’s initialized beforehand.
- The first choice returned by `GenerateContent` is used – if multiple choices are expected, extend the logic to iterate over all.

---

## Summary of what the package does

* Wraps an OpenAI LLM and a tool executor into a single agent.  
* Provides two convenient entry points (`Run`, `SimpleRun`) that automatically inject the configured tools into the LLM call, execute any returned tool calls, and return the AI’s first response.  
* Keeps all configuration in one place (the struct fields), making it easy to swap out the underlying LLM or executor.

---