# util

## Overview
`pkg/util/util.go` implements a small helper that removes an XML‑style `<think>` block from the beginning of a string and trims surrounding whitespace.  
The package exposes one public function, `RemoveThinkTag`, which can be used by other parts of the application or directly from a command‑line entry point.

---

## File structure

```
pkg/
└── util/
    └── util.go
```

---

## Environment variables / configuration options
No explicit environment variables are defined in this file.  
If the package is used as part of a larger build, the following items can be configured:

| Variable | Purpose |
|----------|---------|
| `GOFLAGS` | Build flags for the Go compiler (e.g., `-race`, `-v`). |
| `GOBIN`   | Output directory for compiled binaries. |

---

## Command‑line arguments / flags
The package itself does not expose any command‑line flags, but it can be invoked from a main program that passes an input string to `util.RemoveThinkTag`.  
Typical usage:

```bash
go run ./cmd/main.go --input="…"
```

or

```bash
go build -o bin/util
./bin/util <input>
```

---

## Code walk‑through

### 1. Regular expression definition
```go
var thinkRegex = regexp.MustCompile(`(?s)^<think>.*?</think>`)
```
* Compiles a regex that matches any `<think>` block from the start of a string up to the closing `</think>` tag, including newlines (`(?s)` flag).  
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

## Edge cases & launch scenarios
* **Empty input** – The function returns the original string unchanged if `<think>` is not present.  
* **Multiple `<think>` blocks** – Only the first block at the beginning of the string is removed; subsequent blocks are left untouched.  
* **CLI usage** – If this package is imported by a `main` package, it can be called as `util.RemoveThinkTag(os.Args[1])` or via a flag parser.

---