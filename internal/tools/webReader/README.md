## webreader Package Summary

**Package Name:** `webreader`

This package fetches content from URLs, optionally rendering JavaScript-heavy pages using a headless browser (Rod), and converts the HTML to Markdown format using `html-to-markdown`. It's designed for extracting readable text from webpages.

**Configuration/Arguments:**

*   The primary input is a URL string (`u`) passed to `ProcessUrl`.
*   JavaScript rendering is determined dynamically by checking for `<noscript>` tags containing "javascript" in the HTML source using `checkNoScript`. If present, Rod (headless browser) is used. Otherwise, standard HTTP GET requests are made.

**File Structure:**

```
internal/tools/webReader/
├── webReader.go
```

**Key Functions:**

*   `ProcessUrl(u string)`: Main function to fetch and convert a URL to Markdown. Handles both static HTML and JavaScript-rendered pages.
*   `loadJSPage(u string)`: Launches Rod, navigates to the URL, waits for rendering, and extracts the final HTML content.  Error handling is included for browser connection, page loading, and DOM stabilization.
*   `checkNoScript(url string)`: Determines if JavaScript rendering is needed by parsing the initial HTML source for `<noscript>` tags containing "javascript".

**Dependencies:**

*   `github.com/JohannesKaufmann/html-to-markdown/v2`: For converting HTML to Markdown.
*   `github.com/PuerkitoBio/goquery`: For parsing HTML (used in `checkNoScript`).
*   `github.com/go-rod/rod`: For headless browser automation (JavaScript rendering).

**Edge Cases:**

*   If Rod fails to launch or connect, the function will return an error.  The code includes basic error handling but may not cover all possible failure scenarios with `Rod`.
*   Pages that dynamically load content *after* initial DOM stabilization might not be fully rendered by Rod if the wait time is insufficient.