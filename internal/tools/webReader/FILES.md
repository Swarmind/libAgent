# internal/tools/webReader/webReader.go  
## webreader Package Summary  
  
**Package Name:** `webreader`  
  
**Imports:**  
  
*   `fmt`: For formatted I/O.  
*   `io`: For basic I/O interfaces.  
*   `net/http`: For making HTTP requests.  
*   `strings`: For string manipulation.  
*   `time`: For time-related functions.  
*   `github.com/JohannesKaufmann/html-to-markdown/v2`: HTML to Markdown conversion library.  
*   `github.com/JohannesKaufmann/html-to-markdown/v2/converter`: Converter options for `html-to-markdown`.  
*   `github.com/PuerkitoBio/goquery`: For parsing HTML documents.  
*   `github.com/go-rod/rod`: For browser automation (rendering JavaScript pages).  
*   `github.com/go-rod/rod/lib/proto`: Rod's protocol definitions.  
  
**External Data Sources:**  
  
*   URLs provided as input to `ProcessUrl`. The package fetches content from these URLs using HTTP GET requests or, if necessary, renders them with a headless browser (Rod) for JavaScript-heavy pages.  
*   The HTML content fetched from the URL is processed by `html-to-markdown` to convert it into Markdown format.  
  
**TODOs:** None found in this file.  
  
### Code Sections Summary:  
  
1.  **`ProcessUrl(u string)` Function:** This function is the main entry point for processing a given URL (`u`). It first checks if JavaScript rendering is required using `checkNoScript`. If so, it loads the page with Rod; otherwise, it fetches the HTML directly via HTTP GET. The fetched content (either from HTTP or Rod) is then converted to Markdown using `htmltomarkdown` and returned as a string.  
  
2.  **`loadJSPage(u string)` Function:** This function uses the `go-rod` library to launch a headless browser, navigate to the given URL (`u`), wait for the page to load (including JavaScript execution), and extract the rendered HTML content. It handles potential errors during browser connection, page creation, loading, or DOM stabilization.  
  
3.  **`checkNoScript(url string)` Function:** This function checks if a webpage requires JavaScript rendering by parsing its HTML using `goquery`. It searches for `<noscript>` tags containing the word "javascript" (case-insensitive). If found, it assumes that JavaScript is necessary to render the page correctly and returns true; otherwise, it returns false.  
  
