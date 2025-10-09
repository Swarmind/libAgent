# examples/simple/main.go  
**Package Name:** `main`  
  
**Imports:**  
- `context`  
- `fmt`  
- `os`  
- `github.com/Swarmind/libagent/pkg/agent/simple`  
- `github.com/Swarmind/libagent/pkg/config`  
- `github.com/Swarmind/libagent/pkg/util`  
- `github.com/rs/zerolog`  
- `github.com/rs/zerolog/log`  
- `github.com/tmc/langchaingo/llms/openai`  
  
**External Data / Input Sources:**  
- Configuration loaded from `config.NewConfig()`. This likely reads environment variables or a config file to determine AI URL, token, and model name.  
- The prompt string is defined as a constant: `Prompt = \`This is a test. Write OK in response.\``  
  
**TODOs:** None found in the provided code snippet.  
  
---  
  
### Initialization & Configuration  
  
The code initializes logging with debug level output to stderr using zerolog. It then loads configuration from an unspecified source (likely environment variables or config file) via `config.NewConfig()`. The loaded configuration is used later for OpenAI API setup.  
  
### LLM Setup  
  
An OpenAI language model (`llm`) is initialized using the configured AI URL, token, and model name. Error handling is present to log fatal errors if initialization fails.  The OpenAI version is hardcoded as "v1".  
  
### Agent Execution & Output  
  
A `simple.Agent` instance is created and assigned the initialized LLM. The agent then executes a simple run with a predefined prompt (`Prompt`) using default call options derived from the configuration. Finally, the result (after removing any potential "Think" tags) is printed to standard output.  The context passed into `SimpleRun` is background context.  
  
