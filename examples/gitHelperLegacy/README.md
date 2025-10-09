## Package `main` Summary

This package provides a command-line tool that leverages ReWOO (a large language model agent) to automate GitHub issue resolution. It takes a repository URL, an issue identifier (number, URL, or description), and optional additional context as input. The core logic involves constructing a detailed task description for the ReWOO agent, including instructions on cloning the repository, identifying relevant files, making changes using `git` commands executed via a custom tool (`CommandExecutor`), and verifying those changes before committing them.

### Environment Variables & Flags:

*   `-repo`: Repository URL (default: `https://github.com/JackBekket/Reflexia`).
*   `-issue`: Issue identifier (required). Can be an issue number, full URL, or descriptive text.
*   `-additionalContext`: Optional additional context for the ReWOO agent's task.

### File Structure:

*   `main.go`: Contains the main application logic, including argument parsing, tool setup, task description construction, and ReWOO execution.

### Major Code Parts Summary:

1.  **Initialization:** Sets up logging with `zerolog` at debug level and loads configuration using `config.NewConfig()`.
2.  **Argument Parsing:** Parses command-line flags for repository URL, issue input, and additional context. Validates that the `-issue` flag is provided.
3.  **Tool Setup:** Creates a `toolsExecutor` instance with whitelisted tools: `ReWOOTool`, `CommandExecutor`, `SemanticSearchTool`, and `WebReaderTool`. The executor handles cleanup using `defer`.
4.  **Repository Details Extraction:** Extracts the repository name from the URL, defaulting to "Reflexia" if the default repo is used.
5.  **Task Description Construction:** Constructs a detailed task description for ReWOO, including instructions on cloning the repository, finding and modifying files using `CommandExecutor` (which executes shell commands), verifying changes with `git status`, and handling multi-line modifications. The code warns about exceeding 4000 character limits for task descriptions.
6.  **ReWOO Execution:** Marshals ReWOO query arguments into JSON, calls the ReWOO tool using the tools executor, prints the result to stdout, and provides operational reminders (git/gh installation, authentication, pgvector database setup).

The code is designed as a self-contained example of how to integrate ReWOO with external tools for automating GitHub issue resolution. The task description is highly detailed, providing step-by-step instructions for the agent to follow. Error handling and logging are implemented throughout the process.  The `CommandExecutor` tool allows arbitrary shell commands to be executed within the workflow, which could pose security risks if not carefully controlled.