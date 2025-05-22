# tools

## Overview
The `tools` package provides a collection of integrated tools for system operations, including web search, command execution, network scanning, semantic search, and more. Each tool is configurable via `config.Config` and designed to work within a unified execution framework.

## Package Structure
The package contains the following files, each implementing a distinct tool:

- `ddgSearch.go`: DDG Search (DuckDuckGo) integration
- `executor.go`: Command execution in temporary directories
- `nmap.go`: Nmap network scanning
- `rewoo.go`: ReWOO complex task handling
- `semanticSearch.go`: Semantic code search via vector databases
- `tools.go`: Core tool management and registration
- `webReader.go`: Web content retrieval via URLs

## Configuration Options
The package respects the following configuration values in `config.Config`:

### Tool Control
- `DDGSearchDisable` (bool): Disable DDG search
- `NmapDisable` (bool): Disable Nmap
- `ReWOODisable` (bool): Disable ReWOO
- `SemanticSearchDisable` (bool): Disable semantic search
- `WebReaderDisable` (bool): Disable web reader

### Tool Configuration
- `DDGSearchMaxResults` (int): Max results for DDG searches (default: 5)
- `DDGSearchUserAgent` (string): Custom User-Agent for DDG
- `CommandExecutorDisable` (bool): Disable command executor
- `CommandExecutorCommands` (map[string]string): Command descriptions
- `NmapDisable` (bool): Disable Nmap
- `SemanticSearchAIURL` (string): OpenAI API endpoint
- `SemanticSearchAIToken` (string): OpenAI API token
- `SemanticSearchDBConnection` (string): Database connection string
- `SemanticSearchEmbeddingModel` (string): Embedding model to use
- `SemanticSearchMaxResults` (int): Max search results (default: 2)
- `WebReaderDisable` (bool): Disable web reader

## Notes
- **Unimplemented Features**: 
  - `executor.go`: `sessionDir` is declared but unused
  - `rewoo.go`: Typo `config.ConifgToCallOptions` (should be `ConfigToCallOptions`)
  - `semanticSearch.go`: `collection` parameter may conflict with function name
- **Edge Cases**: 
  - Zero `MaxResults` is treated as 2 in semantic search
  - Nmap fails silently if not installed
  - Invalid JSON input to `Call` methods returns errors
- **Tool Registration**: All tools are registered in `globalToolsRegistry` during initialization.