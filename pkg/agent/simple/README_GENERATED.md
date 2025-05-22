# simple

The `simple` package provides an agent that uses an OpenAI LLM to generate responses. It includes two methods: `Run` for handling chat states and `SimpleRun` for single-message interactions.

**External Dependencies**  
- The `openai.LLM` instance must be configured externally (e.g., API keys, model parameters). No configuration options are exposed by this package.

**Limitations**  
- No handling for empty chat state in `Run` method.  
- No error handling for empty response `Choices[0].Content`.  
- No validation for `LLM` instance being `nil`.  
- No rate-limiting or retry logic implemented.  

**Notes**  
- The package assumes `LLM.GenerateContent` returns at least one choice.  
- No TODOs or comments present in the code.