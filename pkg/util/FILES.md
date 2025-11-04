# pkg/util/util.go  
**Package / Component**    
`util`  
  
---  
  
### Imports  
| Package | Purpose |  
|---------|---------|  
| `regexp` | Compile a regular expression used to locate `<think>` tags in the input string. |  
| `strings` | Provide string manipulation utilities (`HasPrefix`, `TrimSpace`). |  
  
---  
  
### External Data / Input Sources    
* The function operates on an arbitrary string passed as its argument; no external files or services are referenced.  
  
---  
  
### TODOs  
No explicit TODO comments were found in this file.  
  
---  
  
## Summary of Major Code Parts  
  
### 1. Regular Expression Definition  
```go  
var thinkRegex = regexp.MustCompile(`(?s)^<think>.*?</think>`)  
```  
* Compiles a regex that matches any `<think>` block from the beginning of a string up to the closing `</think>` tag, including newlines (`(?s)` flag).    
* The compiled pattern is stored in a package‑level variable for reuse by other functions.  
  
### 2. Function `RemoveThinkTag`  
```go  
func RemoveThinkTag(input string) string {  
    if !strings.HasPrefix(input, "<think>") {  
        return input  
    }  
  
    return strings.TrimSpace(thinkRegex.ReplaceAllString(input, ""))  
}  
```  
* **Purpose**: Strip the leading `<think>` block from an input string and trim surrounding whitespace.    
* **Logic**:  
  * Checks whether the input starts with `<think>`; if not, it returns the original string unchanged.  
  * Uses `thinkRegex.ReplaceAllString` to replace the matched block with an empty string, effectively removing it.  
  * Wraps the result in `strings.TrimSpace` to eliminate any leading/trailing whitespace that may remain after removal.  
  
---  
  
**<end_of_output>**  
  
