# Package `main` – Git Helper Legacy

## Project structure
```
examples/gitHelperLegacy/
└── main.go
```

---

## Imports & Purpose  

| Import | Purpose |
|--------|---------|
| `context` | Provides background context for tool execution |
| `encoding/json` | JSON marshaling of ReWOO query arguments |
| `flag` | Command‑line flag parsing (`repo`, `issue`, `additionalContext`) |
| `fmt` | String formatting (task description, logs) |
| `os` | Access to standard error output for logger |
| `path/filepath` | Path manipulation for repo name extraction |
| `strings` | String trimming & splitting |
| `github.com/Swarmind/libagent/pkg/config` | Application configuration loader |
| `github.com/Swarmind/libagent/pkg/tools` | Tool definitions and executor |
| `github.com/rs/zerolog` | Logging level control |
| `github.com/rs/zerolog/log` | Logger instance |

---

## Environment variables, flags & command‑line arguments  

* **Flags**  
  * `-repo` – defaults to the constant `defaultRepoURL`.  
  * `-issue` – required; can be a number, full URL, or descriptive text.  
  * `-additionalContext` – optional string that is appended to the ReWOO prompt.

* **Environment variables** – none explicitly referenced in this file.

* **Command‑line usage**  
  ```bash
  go run examples/gitHelperLegacy/main.go -repo <url> -issue <id|url|text> [-additionalContext <txt>]
  ```

---

## How the application can be launched (edge cases)

1. **Default launch** – no flags supplied; uses hard‑coded defaults for `repo` and empty context.  
2. **Custom repo URL** – supply a full GitHub repository URL via `-repo`.  
3. **Issue as number** – pass an issue number (`-issue 42`) which will be interpreted by the ReWOO tool.  
4. **Full context string** – add extra context with `-additionalContext` to influence the prompt.

---

## Code logic summary

### 1. Logging & flag setup
```go
zerolog.SetGlobalLevel(zerolog.DebugLevel)
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
```
Initializes global logger and sets output to `stderr`. Flags are defined, parsed, and validated (`issue` must be non‑empty).

### 2. Configuration & context
```go
cfg, err := config.NewConfig()
ctx := context.Background()
```
Loads application configuration and creates a background context for tool execution.

### 3. Tools whitelist & executor creation
```go
toolsToWhitelist := []string{ ... }
toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))
```
Specifies which tools the ReWOO agent may use and creates an executor instance.

### 4. Repository details derivation  
The code trims a trailing `.git` from the supplied URL, then:
* If it matches the default repo URL → sets `repoName = "Reflexia"` and owner/repo string accordingly.
* Otherwise derives the base name and splits the path to build an `ownerAndRepo` string.

### 5. Task description construction  
A large multi‑line template is built with `fmt.Sprintf`.  
It contains a step‑by‑step plan for the ReWOO agent, including:
1. Issue understanding (WebReader → LLM → SemanticSearch).  
2. Development environment setup (git clone, cd, safe.directory config, tree/ls).  
3. File discovery & modification steps with placeholders for file paths and commands.  
4. Verification steps (cat, git status, diff).

The template is populated with the parsed flags and derived repo name.

### 6. ReWOO query execution  
```go
rewooQueryArgs := tools.ReWOOToolArgs{ Query: taskDescription }
rewooQueryBytes, err := json.Marshal(rewooQueryArgs)
result, err := toolsExecutor.CallTool(ctx,
    tools.ReWOOToolDefinition.Name,
    string(rewooQueryBytes),
)
```
The constructed description is marshaled to JSON and sent to the ReWOO tool via the executor. Result output is printed to stdout.

### 7. Logging & output
Informational logs record repository, issue, context length, and task description size.  
After execution, a summary of the agent’s report and operational reminders are printed.

---

## Relations between code entities

* `toolsExecutor` is created once at program start; it receives the configuration and a whitelist of tool names.
* The `repoName`, `ownerAndRepo`, and `issue` variables feed directly into the prompt template, ensuring that the ReWOO agent knows which repository to target and what issue to address.
* The JSON marshaling step (`rewooQueryArgs`) is the bridge between Go code and the ReWOO tool; it serializes the prompt string for consumption by the executor.

---

## Unclear places / dead code

No explicit `TODO:` comments or unused imports are present in this file. All referenced variables and functions appear to be used, so no dead code was detected.

---