# pkg/tools/ddgSearch.go  
## Package/Component Name: `tools`  
  
**Imports:**  
  
*   `context`: For managing context in function calls.  
*   `encoding/json`: For JSON serialization and deserialization.  
*   `github.com/Swarmind/libagent/internal/tools`: Internal tools package (likely for tool registration).  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration management within the libagent framework.  
*   `github.com/tmc/langchaingo/llms`: Langchain LLM definitions and interfaces.  
*   `github.com/tmc/langchaingo/tools/duckduckgo`: DuckDuckGo search tool implementation from Langchain.  
  
**External Data / Input Sources:**  
  
*   Configuration (`config.Config`) is used to determine whether the DDG search tool should be enabled, maximum results count and user agent string.  
*   Input `string` (JSON formatted) for the `Call` method of `DDGSearchTool`, which represents the search query.  
  
**TODOs:** None found in this file.  
  
---  
  
### Function Definition: `DDGSearchDefinition`  
  
Defines a Langchain function definition named "webSearch" with a description and parameter schema (query string). This is used to expose the DuckDuckGo search functionality as an LLM tool.  
  
### Struct: `DDGSearchArgs`  
  
A struct that represents the arguments for the DDG search, containing only a single field: `Query` of type `string`. Used for unmarshaling JSON input into structured data.  
  
### Struct: `DDGSearchTool`  
  
Wraps a DuckDuckGo tool (`duckduckgo.Tool`) and provides a `Call` method to execute searches based on the provided query string. The wrapped tool is initialized in the `init()` function.  
  
### Function: `DDGSearchTool.Call`  
  
Unmarshals JSON input into a `DDGSearchArgs` struct, then calls the underlying DuckDuckGo search tool with the extracted query. Returns the search results as a string or an error if unmarshalling fails.  
  
### Function: `init()`  
  
Registers the DDG search tool within a global tools registry (`globalToolsRegistry`). The registration is conditional based on configuration settings (disabled flag, max results, user agent). If enabled, it creates and wraps a DuckDuckGo search instance using configured parameters before registering it as a callable tool.  Handles default values for `DDGSearchMaxResults` and `DDGSearchUserAgent`.  
  
# pkg/tools/executor.go  
**Package/Component Name:** `tools`  
  
**Imports:**  
  
*   `context`  
*   `encoding/json`  
*   `fmt`  
*   `os`  
*   `strings`  
*   `time`  
*   `github.com/Swarmind/libagent/internal/tools`  
*   `github.com/Swarmind/libagent/pkg/config`  
*   `github.com/ThomasRooney/gexpect`  
*   `github.com/google/uuid`  
*   `github.com/rs/zerolog/log`  
*   `github.com/tmc/langchaingo/llms`  
  
**External Data/Input Sources:**  
  
*   Configuration (`config.Config`) to disable the command executor or define allowed commands with descriptions.  
*   JSON input string for `CommandExecutorArgs`.  
*   Environment variables (specifically, `PATH`).  
  
**TODOs:** None found in this code snippet.  
  
---  
  
### Command Executor Definition and Arguments  
  
The code defines a LangChain-compatible function definition (`CommandExecutorDefinition`) for executing shell commands. It includes the command as input via JSON payload. The `CommandExecutorArgs` struct is used to parse the incoming JSON data, expecting a "command" field.  
  
### Command Execution Logic  
  
The core functionality resides in the `RunCommand` method of the `CommandExecutorTool`. This tool creates a temporary directory for session isolation and spawns an interactive bash shell using `gexpect`. It sets up a unique prompt (UUID) to reliably detect command completion, captures output, and executes commands. The code handles potential errors during spawning or execution by sending Ctrl+C (`0x03`) if the expected prompt isn't received within 30 seconds.  
  
### Temporary Directory Management  
  
The `cleanup` method removes the temporary directory created for each session to prevent resource leaks. It also closes the spawned process. The tool creates a temp dir only when it is not already initialized, and reuses existing one in subsequent calls.  
  
### Tool Registration with Configuration  
  
The `init` function registers the command executor as part of a global tools registry (`globalToolsRegistry`). This registration is conditional based on configuration settings (`cfg.CommandExecutorDisable`). If enabled, the tool's description can be augmented with custom commands and their descriptions from the config file. The code ensures that trailing newlines are removed from the final function definition to avoid formatting issues.  
  
# pkg/tools/metasploit.go  
## Package/Component Summary: `tools` - Metasploit Search Tool  
  
**Package Name:** `tools`  
  
**Imports:**  
  
*   `context`: For managing execution context.  
*   `encoding/json`: For JSON serialization and deserialization.  
*   `fmt`: For formatted I/O operations.  
*   `os/exec`: For executing external commands (Metasploit).  
*   `strings`: For string manipulation.  
*   `github.com/Swarmind/libagent/internal/tools`: Internal tools package dependency.  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration package dependency.  
*   `github.com/tmc/langchaingo/llms`: Langchain LLM definitions for tool integration.  
  
**External Data/Input Sources:**  
  
*   Metasploit executable path (configurable via `SetMsfCommand`). Defaults to "msfconsole".  
*   JSON input containing a list of Metasploit search queries (`Queries` field).  
*   Configuration settings from `config.Config`, specifically `cfg.MsfDisable` which can disable the tool.  
*   Port information (from `GenerateMsfQueries`) used to construct default exploit searches based on open ports and services.  
  
**TODOs:** None found in this file.  
  
---  
  
### Code Sections Summary:  
  
1.  **Function Definitions (`MsfSearchToolDefinition`):** Defines the Langchain function definition for the Metasploit search tool, including its name, description, and input parameters (a list of queries). The `queries` parameter accepts a list of strings that will be used as Metasploit search terms.  
  
2.  **MsfSearchTool Struct & Methods:** Implements the core logic for executing Metasploit searches.  
    *   The `MsfSearchTool` struct holds configurable values for the executable path and argument template.  
    *   The `Call` method takes a JSON string as input, unmarshals it into a `MsfSearchToolArgs` struct (containing the list of queries), executes Metasploit searches using `exec.Command`, captures the output, and returns the results in JSON format.  
  
3.  **Query Generation (`GenerateMsfQueries`):** Generates default Metasploit search queries based on a list of open ports and their associated services. It creates exploit-focused queries like `"type:exploit name:<service>"` and port-specific searches.  
  
4.  **Initialization & Registration (`init()`):** Registers the `MsfSearchTool` with a global tools registry (presumably within the larger application). The tool is only registered if `cfg.MsfDisable` is false, allowing for easy disabling of Metasploit functionality via configuration.  
  
5. **Configuration:** Allows to set custom executable and arguments template using `SetMsfCommand`.  
  
# pkg/tools/nmap.go  
**Package/Component Name:** `tools`  
  
**Imports:**  
*   `context`  
*   `encoding/json`  
*   `fmt`  
*   `net`  
*   `os/exec`  
*   `regexp`  
*   `strings`  
*   `github.com/Swarmind/libagent/internal/tools`  
*   `github.com/Swarmind/libagent/pkg/config`  
*   `github.com/tmc/langchaingo/llms`  
  
**External Data/Input Sources:**  
*   JSON input string for `NmapToolArgs` (IP address and optional arguments).  
*   Configuration (`config.Config`) to disable the tool via `cfg.NmapDisable`.  
  
**TODOs:**  
*   "should be in another pkg" - Refers to the `GenerateMsfQueries` function, suggesting it should be moved to a separate package for better organization.  
  
---  
  
### Nmap Tool Definition and Arguments Handling  
The code defines an LLM function definition (`NmapToolDefinition`) for executing nmap scans with configurable arguments. The input is expected as JSON containing the target IP address and optional nmap arguments. If no arguments are provided, default values (verbosity, timing template, TCP connect scan, version detection, host discovery, full version info, top ports) are used. Input validation checks if the provided IP address is valid before proceeding.  
  
### Nmap Execution  
The `Call` function executes the nmap command using `exec.Command`. The output of the command is captured and parsed to extract port information. Error handling includes logging the nmap output in case of failure.  
  
### Port Parsing (`ParseNmapPorts`)  
This function uses a regular expression to parse the raw nmap output, extracting port numbers, states (open, closed, filtered), and service names. The extracted data is stored in a `PortInfo` struct and returned as a slice. Debug print statements are included for logging found ports, state, and service information.  
  
### Metasploit Query Generation (`GenerateMsfQueries`)  
The code generates Metasploit search queries based on the parsed port information (TODO: should be moved to another package). The exact implementation of this function is not provided in the snippet but it's assumed to use the extracted ports to create relevant msfconsole commands.  
  
### Tool Registration (`init`)  
The `init` function registers the Nmap tool with a global registry (`globalToolsRegistry`). This allows other parts of the system to discover and call the nmap functionality. The registration is conditional based on a configuration flag (`cfg.NmapDisable`), allowing the tool to be disabled if necessary.  
  
# pkg/tools/rewoo.go  
## Package/Component Name: `tools`  
  
**Imports:**  
  
*   `context`: For managing context in function calls.  
*   `encoding/json`: For JSON serialization and deserialization.  
*   `github.com/Swarmind/libagent/internal/tools`: Internal tools package dependency.  
*   `github.com/Swarmind/libagent/internal/tools/rewoo`: Specific rewoo-related internal tool functionality.  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration handling from the main package.  
*   `github.com/JackBekket/langgraphgo/graph/stategraph`: LangGraph state graph implementation.  
*   `github.com/tmc/langchaingo/llms`: LangChain LLM interface.  
*   `github.com/tmc/langchaingo/llms/openai`: OpenAI LLM integration for LangChain.  
  
**External Data / Input Sources:**  
  
*   Configuration (`config.Config`) loaded from the main package, including AI URL, token, model name, and ReWOO-specific settings (e.g., `ReWOODisable`, `RewOODefaultCallOptions`).  
*   Input string to the `Call` method is expected in JSON format conforming to `ReWOOToolArgs`.  
  
**TODOs:** None found in this file.  
  
---  
  
### Code Summary:  
  
**1. ReWOO Tool Definition (`ReWOOToolDefinition`)**: Defines a LangChain-compatible function definition for the "rewoo" tool, including its name, description (focused on complex reasoning tasks), and input parameter schema ("query" as a string). This is used to expose the tool in an LLM context.  
  
**2. `ReWOOToolArgs` Struct**: Defines the expected JSON structure for input arguments passed to the ReWOO tool via the `Call` method, containing only a "query" field (string).  
  
**3. `ReWOOTool` Struct**: Represents the ReWOO tool itself. It holds an instance of `rewoo.ReWOO` (presumably handling core rewoo logic) and a LangGraph runnable graph (`graph`). The graph is initialized lazily when needed.  
  
**4. `Call` Method**: This method handles incoming requests to execute the ReWOO tool:  
    *   Parses JSON input into `ReWOOToolArgs`.  
    *   Ensures that the internal tools executor of `rewoo.ReWOO` is initialized (using a global registry).  
    *   Initializes the LangGraph state graph if it hasn't been already.  
    *   Invokes the graph with the provided query as input, retrieves the result from the final state, and returns it.  
  
**5. `init` Function**: Registers the ReWOO tool within a global tools registry (`globalToolsRegistry`). This function is executed when the package is initialized:  
    *   Checks if the ReWOO tool is disabled via configuration (`cfg.ReWOODisable`). If so, registration is skipped.  
    *   Creates an OpenAI LLM instance using configured parameters (AI URL, token, model).  
    *   Instantiates `ReWOOTool` with the created LLM and stores it in a global registry for later use.  
  
# pkg/tools/semanticSearch.go  
```text  
Package Name: `tools`  
  
Imports:  
- `context`  
- `encoding/json`  
- `fmt`  
- `github.com/Swarmind/libagent/internal/tools`  
- `github.com/Swarmind/libagent/pkg/config`  
- `github.com/jackc/pgx/v5/pgxpool`  
- `github.com/tmc/langchaingo/embeddings`  
- `github.com/tmc/langchaingo/llms`  
- `github.com/tmc/langchaingo/llms/openai`  
- `github.com/tmc/langchaingo/vectorstores/pgvector`  
  
External Data / Input Sources:  
- OpenAI URL (`cfg.SemanticSearchAIURL`)  
- OpenAI Token (`cfg.SemanticSearchAIToken`)  
- PostgreSQL Connection String (`cfg.SemanticSearchDBConnection`)  
- Embedding Model Name (`cfg.SemanticSearchEmbeddingModel`)  
- Max Results for Search (`cfg.SemanticSearchMaxResults`, defaults to 2 if not provided)  
- Input string (JSON formatted `SemanticSearchArgs` containing `query` and `collection`).  
  
TODOs:  
- "there should NOT exist arguments which called NAME cause it cause COLLISION with actual function name.    .....more like confusion then collision so there are no error" - This comment suggests a potential naming conflict in the argument definitions, but doesn't specify an immediate fix.  
  
Code Summary:  
  
### Semantic Search Definition  
The `SemanticSearchDefinition` is a Langchain FunctionDefinition that describes the semantic search tool. It takes a query string and collection name as input and returns matching file contents from a vector store. The definition includes parameter validation for both inputs, ensuring they are strings with appropriate descriptions.  
  
### Semantic Search Tool Implementation (`SemanticSearchTool`)  
The `SemanticSearchTool` struct holds configuration parameters (OpenAI URL, token, DB connection, embedding model, max results). Its `Call` method performs the actual semantic search:  
1.  Parses JSON input into a `SemanticSearchArgs` struct.  
2.  Establishes a PostgreSQL connection using `pgxpool`.  
3.  Initializes an OpenAI LLM and Embedder.  
4.  Creates a `pgvector` store with the specified collection name, connecting to the database and embedding model.  
5.  Performs similarity search using the query string and max results.  
6.  Concatenates matching page contents into a response string.  
  
### Tool Registration (`init`)  
The `init` function registers the semantic search tool within a global registry (`globalToolsRegistry`). It checks for required configuration parameters (OpenAI URL, token, DB connection, embedding model) and returns an error if any are missing. If all parameters are present, it creates a `SemanticSearchTool` instance with the provided values and adds it to the registry as a callable tool. The function also disables the tool if `cfg.SemanticSearchDisable` is true.  
  
### Data Structures  
-   `SemanticSearchArgs`: Struct for parsing JSON input containing query and collection name.  
-   `SemanticSearchTool`: Holds configuration parameters for OpenAI, DB connection, embedding model, and max results.  
<end_of_output>  
```  
  
# pkg/tools/tools.go  
## Package: `tools`  
  
**Imports:**  
  
*   `context`: For managing context in asynchronous operations.  
*   `slices`: Provides slice manipulation functions (e.g., `Contains`).  
*   `github.com/Swarmind/libagent/internal/tools`: Internal tools package, likely containing tool definitions and execution logic.  
*   `github.com/Swarmind/libagent/pkg/config`: Package for configuration management.  
  
**External Data / Input Sources:**  
  
*   `context.Context`: Used as input to initialize tools.  
*   `config.Config`: Configuration object passed during tool initialization.  
*   `globalToolsRegistry`: A slice of functions that initialize individual tools, taking context and config as arguments. This is the primary mechanism for registering available tools.  
*   `ExecutorOption`: Function type used to configure `ExecutorOptions`.  
  
**TODOs:** None found in this file.  
  
---  
  
### Tool Registration & Initialization  
  
The core functionality revolves around the `globalToolsRegistry`, which holds functions responsible for initializing individual tools. The `NewToolsExecutor` function iterates through this registry, calling each initialization function with a provided context and configuration. Tools are added to an internal map (`tools`) if they initialize successfully and pass any whitelist filtering (if configured).  
  
### Executor Options & Whitelisting  
  
The `ExecutorOptions` struct allows configuring the tool executor, specifically via a `ToolsWhitelist`. The `WithToolsWhitelist` function provides a convenient way to set this whitelist. Tools not present in the whitelist are skipped during initialization. This mechanism enables selective loading of tools based on configuration or runtime requirements.  
  
### Global Executor Instance  
  
The code maintains a global instance of `tools.ToolsExecutor` (`globalToolsExecutor`). The initialized executor is assigned to this global variable, suggesting that it's intended as a singleton for the entire application lifecycle.  
  
---  
  
# pkg/tools/webReader.go  
## Package/Component Name and Imports  
  
**Package Name:** `tools`  
  
**Imports:**  
  
*   `context`: For managing context in function calls.  
*   `encoding/json`: For JSON parsing of input arguments.  
*   `github.com/Swarmind/libagent/internal/tools`: Internal tools package dependency.  
*   `github.com/Swarmind/libagent/internal/tools/webReader`: Specific web reader implementation within internal tools.  
*   `github.com/Swarmind/libagent/pkg/config`: Configuration access for disabling the tool.  
*   `github.com/tmc/langchaingo/llms`: Langchain LLM definitions and structures.  
  
## External Data / Input Sources  
  
The code accepts a URL as input via JSON string, which is then processed by `webreader.ProcessUrl`. The configuration (`cfg`) from the `config.Config` type determines whether this tool is enabled or disabled.  
  
## TODOs  
  
There are no explicit `TODO` comments in the provided code snippet.  
  
---  
  
### WebReader Definition and Arguments  
  
The `WebReaderDefinition` defines a Langchain function with a single parameter: "url" (string). The description emphasizes that only valid URLs should be passed, suggesting LLM-based URL extraction as a prerequisite step.  `WebReaderArgs` is a struct to hold the parsed JSON input containing the URL.  
  
### WebReader Tool Implementation  
  
The `WebReaderTool` struct contains the core logic for processing the URL. The `Call` method unmarshals the input string (expected in JSON format) into a `WebReaderArgs` instance and then calls `webreader.ProcessUrl` to extract text from the provided URL. Error handling is included during JSON parsing.  
  
### Tool Registration  
  
The `init()` function registers this tool with a global registry (`globalToolsRegistry`). The registration process checks for a configuration flag (`cfg.WebReaderDisable`) which, if true, prevents the tool from being registered. If enabled, an instance of `WebReaderTool` is created and its definition and call method are added to the registry as a `tools.ToolData`.  
  
