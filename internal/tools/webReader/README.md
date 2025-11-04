# webreader  

A lightweight helper that pulls a web page (optionally rendered with JavaScript) and turns it into Markdown.  
It lives under **internal/tools/webReader** and is intended to be used by other parts of the project or run directly from a CLI.

---

## File structure

```
internal/
└─ tools/
   └─ webReader/
      └─ webReader.go
```

Only one source file exists, so all logic is contained in `webReader.go`.

---

## Core functions & flow

| Function | Purpose | Key steps |
|----------|---------|-----------|
| **`ProcessUrl(u string)`** | Entry point that decides whether the page needs JS rendering and returns Markdown. | 1️⃣ Calls `checkNoScript`. <br>2️⃣ If a `<noscript>` element with “javascript” is found, it uses Rod to load the page (`loadJSPage`). <br>3️⃣ Otherwise it performs a plain HTTP GET.<br>4️⃣ The resulting body (as an `io.Reader`) is fed into `htmltomarkdown.ConvertReader` with the domain set via `converter.WithDomain(u)`. |
| **`checkNoScript(url string)`** | Detects whether the target page contains a `<noscript>` tag that indicates JavaScript execution is required. | 1️⃣ HTTP GET → body.<br>2️⃣ Build a goquery document and look for `<noscript>` elements containing “javascript”.<br>3️⃣ Return `true/false` + error. |
| **`loadJSPage(u string)`** | Loads the page in a headless browser (Rod) when JS is needed. | 1️⃣ Create & connect a Rod browser.<br>2️⃣ Open a new page for `u`. <br>3️⃣ Wait until the page has fully loaded (`WaitLoad`) and the DOM is stable (`WaitDOMStable` with a 5 s timeout).<br>4️⃣ Return the raw HTML string. |

The three functions are tightly coupled:  
*`ProcessUrl` orchestrates everything, delegating to `checkNoScript` and optionally `loadJSPage`, then converting the fetched content into Markdown.*

---

## Environment variables / flags / command‑line arguments

| Variable / flag | Description | Default / usage |
|------------------|-------------|-----------------|
| `WEBREADER_URL` | URL of the page to fetch. | Passed directly to `ProcessUrl`. |
| `WEBREADER_TIMEOUT` | Optional timeout for Rod’s wait operations (currently hard‑coded 5 s). | Can be overridden by adding a flag or env var if you extend the code. |

If this package is used as a CLI tool, a typical invocation would look like:

```bash
go run internal/tools/webReader/main.go -url https://example.com/page
```

or, when integrated into another binary:

```go
import "github.com/yourrepo/internal/tools/webReader"

func main() {
    md, err := webreader.ProcessUrl(os.Args[1])
    // …
}
```

---

## Edge cases & launch scenarios

| Scenario | What to watch for |
|----------|-------------------|
| **Page needs JS** | `checkNoScript` must correctly detect a `<noscript>` element. If the tag is missing, the function falls back to a simple HTTP GET. |
| **Rod fails** | Errors from Rod (e.g., connection or wait failures) are bubbled up by `ProcessUrl`. Ensure network connectivity and that Chrome/Chromium is available. |
| **Timeout too short** | The 5 s timeout in `loadJSPage` may need adjustment for very slow pages; consider exposing it as a flag. |

---

## Summary

* `webReader.go` implements a small pipeline: fetch → optional JS rendering → Markdown conversion.  
* It relies on Rod (headless browser) and html-to-markdown to handle dynamic content, while goquery is used only for the initial `<noscript>` check.  
* All configuration is currently internal; you can extend it by adding flags or env vars for timeout, user‑agent, etc.

---