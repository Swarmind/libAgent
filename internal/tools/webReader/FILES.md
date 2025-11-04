# internal/tools/webReader/webReader.go  
**Package / Component name**    
`webreader`  
  
---  
  
### Imports  
```go  
import (  
	"fmt"  
	"io"  
	"net/http"  
	"strings"  
	"time"  
  
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"  
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"  
	"github.com/PuerkitoBio/goquery"  
	"github.com/go-rod/rod"  
	"github.com/go-rod/rod/lib/proto"  
)  
```  
  
* `fmt`, `io`, `net/http`, `strings`, `time` – standard library utilities.    
* `htmltomarkdown` & `converter` – convert HTML to Markdown.    
* `goquery` – jQuery‑like DOM traversal for Go.    
* `rod` & `proto` – headless browser automation (Chrome/Chromium).  
  
---  
  
### External data / input sources  
| Function | Purpose | Input | Output |  
|----------|---------|-------|--------|  
| `ProcessUrl(u string)` | Main entry point; decides whether to load page via JS or plain HTTP, then converts it to Markdown. | URL string | Markdown string + error |  
| `loadJSPage(u string)` | Loads a JavaScript‑rendered page using Rod and returns its HTML. | URL string | Page HTML string + error |  
| `checkNoScript(url string)` | Detects if the target page contains a `<noscript>` tag that indicates JS is required. | URL string | Boolean flag + error |  
  
---  
  
### TODO comments  
*None found in this file.*  
  
---  
  
## Summary of major code parts  
  
#### 1. `ProcessUrl`  
- Calls `checkNoScript` to determine whether the page needs JavaScript rendering.  
- If JS is needed, it uses `loadJSPage`; otherwise a simple HTTP GET fetches the body.  
- The resulting `io.Reader` (`body`) is passed to `htmltomarkdown.ConvertReader`, with the domain set via `converter.WithDomain(u)`.  
- Returns the Markdown string and any error that occurred.  
  
#### 2. `loadJSPage`  
- Creates a Rod browser instance, connects it, and opens a new page targeting the given URL.  
- Waits for the page to load fully (`WaitLoad`) and for the DOM to stabilize (`WaitDOMStable` with a 5 s timeout).  
- Returns the raw HTML of the rendered page.  
  
#### 3. `checkNoScript`  
- Performs an HTTP GET on the supplied URL, checks that the response status is OK.  
- Builds a GoQuery document from the response body and searches for `<noscript>` elements.  
- If any such element contains the word “javascript”, it flags that JavaScript rendering is required.  
  
---  
  
All functions are tightly coupled: `ProcessUrl` orchestrates the flow, delegating to `checkNoScript` and optionally `loadJSPage`, then converting the fetched content into Markdown.    
The code relies on Rod for headless browser automation and on html-to-markdown for conversion, making it suitable for scraping pages that may need JS execution before being parsed.  
  
