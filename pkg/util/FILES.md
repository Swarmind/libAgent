# pkg/util/util.go  
## Package: `util`  
  
**Imports:**  
  
*   `regexp`: For regular expression operations.  
*   `strings`: For string manipulation functions.  
  
**External Data/Inputs:**  
  
The function `RemoveThinkTag` takes a single string input (`input`). This string is expected to potentially contain `<think>` tags, which the function aims to remove. The presence of these tags determines whether any modification occurs.  
  
**TODOs:**  
  
None found in this code snippet.  
  
### Code Summary:  
  
This package provides utility functions for manipulating strings, specifically targeting removal of a custom tag (`<think>`).   
  
The core functionality resides within `RemoveThinkTag`. This function checks if the input string starts with `<think>`. If not, it returns the original string unchanged. Otherwise, it uses a pre-compiled regular expression (`thinkRegex`) to remove all occurrences of the `<think>` and `</think>` tags (including any content between them) from the input string. The result is then trimmed of leading/trailing whitespace before being returned.  
  
The regex `(?s)^<think>.*?</think>` matches:  
*   `(?s)`: Enables dotall mode, allowing `.` to match newline characters within the tag.  
*   `^`: Matches the beginning of the string.  
*   `<think>`: Matches the opening tag literally.  
*   `.*?`: Matches any character (except newline) zero or more times, non-greedily.  
*   `</think>`: Matches the closing tag literally.  
  
