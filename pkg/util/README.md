## Package: `util`

**Summary:** This package provides a utility function to remove `<think>` tags from strings using regular expressions. The primary use case appears to be stripping out content enclosed within these custom tags, potentially for processing or cleaning text data.

**Configuration/Arguments:**

*   The main configuration is the input string passed to `RemoveThinkTag`.
*   No environment variables, flags, or command-line arguments are used directly by this package. The regex pattern is hardcoded.

**Edge Cases:**

*   If the input string does not start with `<think>`, it's returned unchanged. This means the function only operates on strings that begin with the tag.
*   The regular expression `(?s)^<think>.*?</think>` assumes the closing tag `</think>` exists and is properly nested within the opening tag. Malformed or missing tags could lead to unexpected behavior (e.g., incomplete removal).

**File Structure:**

```
pkg/util/
  - util.go
```

**Code Entities & Relations:**

*   `RemoveThinkTag`: The core function that performs the tag removal logic. It relies on a precompiled regular expression for pattern matching and replacement.
*   `thinkRegex`: A compiled regular expression used by `RemoveThinkTag`. Its hardcoded pattern dictates how tags are matched and removed.

**Potential Issues/Dead Code:**

The regex is anchored to the beginning of the string (`^`). This means it will only remove `<think>` tags if they appear at the very start of the input. If the tag appears elsewhere, it won't be affected. The code doesn't handle nested or overlapping tags correctly; it assumes a simple structure where each opening tag has a corresponding closing tag.