# examples/gitHelperLegacy/main.go  
# Package & Imports    
**Package name:** `main`    
  
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
  
## External Data & Input Sources    
* **Command‑line flags** – `-repo`, `-issue`, `-additionalContext`.    
  * `repo`: defaults to the constant `defaultRepoURL`.    
  * `issue`: required; can be a number, full URL, or descriptive text.    
  * `additionalContext`: optional string that is appended to the ReWOO prompt.    
  
* **Configuration** – loaded via `config.NewConfig()` (no external file path shown).    
  
* **Tools executor** – created with a whitelist of four tool definitions:    
  - `ReWOOToolDefinition.Name`    
  - `CommandExecutorDefinition.Name`    
  - `SemanticSearchDefinition.Name`    
  - `WebReaderDefinition.Name`  
  
---  
  
## TODO List    
No explicit `TODO:` comments are present in the file.    
  
---  
  
## Summary of Major Code Parts  
  
### 1. Logging & Flag Setup  
```go  
zerolog.SetGlobalLevel(zerolog.DebugLevel)  
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})  
```  
Initializes global logger and sets output to `stderr`.    
Flags are defined, parsed, and validated (`issue` must be non‑empty).    
  
### 2. Configuration & Context  
```go  
cfg, err := config.NewConfig()  
ctx := context.Background()  
```  
Loads application configuration and creates a background context for tool execution.  
  
### 3. Tools Whitelist & Executor Creation  
```go  
toolsToWhitelist := []string{ ... }  
toolsExecutor, err := tools.NewToolsExecutor(ctx, cfg, tools.WithToolsWhitelist(toolsToWhitelist...))  
```  
Specifies which tools the ReWOO agent may use and creates an executor instance.  
  
### 4. Repository Details Derivation    
The code trims a trailing `.git` from the supplied URL, then:  
* If it matches the default repo URL → sets `repoName = "Reflexia"` and owner/repo string accordingly.  
* Otherwise derives the base name and splits the path to build an `ownerAndRepo` string.  
  
### 5. Task Description Construction    
A large multi‑line template is built with `fmt.Sprintf`.    
It contains a step‑by‑step plan for the ReWOO agent, including:  
1. Issue understanding (WebReader → LLM → SemanticSearch).    
2. Development environment setup (git clone, cd, safe.directory config, tree/ls).    
3. File discovery & modification steps with placeholders for file paths and commands.    
4. Verification steps (cat, git status, diff).  
  
The template is populated with the parsed flags and derived repo name.  
  
### 6. ReWOO Query Execution    
```go  
rewooQueryArgs := tools.ReWOOToolArgs{ Query: taskDescription }  
rewooQueryBytes, err := json.Marshal(rewooQueryArgs)  
result, err := toolsExecutor.CallTool(ctx,  
    tools.ReWOOToolDefinition.Name,  
    string(rewooQueryBytes),  
)  
```  
The constructed description is marshaled to JSON and sent to the ReWOO tool via the executor.    
Result output is printed to stdout.  
  
### 7. Logging & Output  
Informational logs record repository, issue, context length, and task description size.    
After execution, a summary of the agent’s report and operational reminders are printed.  
  
---  
  
This file serves as the entry point for orchestrating a ReWOO‑based GitHub issue resolution workflow. It gathers user input, prepares an executor with a tool whitelist, constructs a detailed prompt, runs the ReWOO agent, and prints its output.    
All of this will be aggregated into a higher‑level package summary.  
  
