# Package `simple`

The **`pkg/agent/simple`** directory contains a single Go source file that implements an LLM‑powered chat agent.  
It wraps an OpenAI LLM instance and exposes two convenience methods for generating responses from either a full conversation state or a plain string input.

---

## File structure

```
pkg/
└─ agent/
   └─ simple/
      └─ agent.go
```

* `agent.go` – the only file in this package; it defines the public type and its methods.

---

## Imports

```go
import (
    "context"

    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
)
```

* `context` – standard Go context for request handling.  
* `github.com/tmc/langchaingo/llms` – core LangChain‑Go types (e.g., `MessageContent`, `ChatMessageTypeAI/Human`).  
* `github.com/tmc/langchaingo/llms/openai` – the OpenAI LLM implementation used by the agent.

---

## Public type

```go
type Agent struct {
    LLM *openai.LLM
}
```

The agent holds a pointer to an `openai.LLM`.  
All logic in this package revolves around calling `GenerateContent` on that instance.

---

## Methods

| Method | Signature | Purpose |
|--------|------------|---------|
| `Run` | `(ctx context.Context, state []llms.MessageContent, opts ...llms.CallOption) (llms.MessageContent, error)` | Sends the supplied conversation state to the LLM and returns the first choice as an AI chat message. |
| `SimpleRun` | `(ctx context.Context, input string, opts ...llms.CallOption) (string, error)` | Convenience wrapper that builds a single human‑message from `input`, calls the LLM, and returns the raw output string. |

Both methods call `a.LLM.GenerateContent`.  
The first method expects an already‑built slice of message contents; the second constructs that slice on the fly.

---

## Environment variables / configuration

* The package itself does not expose any explicit environment variable or flag.  
  However, the OpenAI LLM instance (`LLM`) can be configured externally (e.g., API key, model name) before creating an `Agent`.

---

## Launch edge cases

If this package is used as a command‑line tool:

1. **As part of a larger binary** – import `"github.com/tmc/langchaingo/pkg/agent/simple"` and create an `Agent` instance in your main program.
2. **Standalone build** – run `go build ./pkg/agent/simple` to produce a binary that can be invoked with `./simple`.  
   The binary would need a small wrapper (e.g., a `main.go`) that creates the agent, reads input from stdin or flags, and calls either `Run` or `SimpleRun`.

---

## Summary

* **Purpose** – provide a thin wrapper around an OpenAI LLM to generate chat responses.  
* **Key entities** – `Agent`, `Run`, `SimpleRun`.  
* **Data flow** – `Run` → `GenerateContent` → first choice → AI message; `SimpleRun` → build human message → same call → raw string.

The package is intentionally minimal, making it easy to embed in larger LangChain‑Go workflows or to expose as a CLI helper.