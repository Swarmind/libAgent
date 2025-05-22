webreader

This package provides functionality to process URLs, convert HTML content to Markdown, and handle JavaScript-rendered web pages. It includes functions to check if a URL requires JavaScript, load pages with or without JavaScript support, and convert HTML to Markdown while preserving domain information in links.

External Data/Config:
- URLs must be valid HTTP(S) endpoints
- HTML content must comply with standard DOM structures
- JS rendering behavior depends on browser automation (go-rod)
- Markdown conversion preserves original domain via `converter.WithDomain()`

Notes:
- No TODOs or comments in code
- Error handling is present but some errors are not explicitly logged
- Resource cleanup (browser/page close) is handled via `defer`
- Conversion may produce inconsistent results for dynamically generated content
- `checkNoScript()` only detects explicit "javascript" references in `<noscript>` tags, not script tags or JS functionality
- `loadJSPage()` has a 5s timeout for DOM stability on complex pages