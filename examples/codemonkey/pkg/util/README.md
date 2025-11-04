<end_of_output>  

# Package `utility`

## Short Summary  
The **`utility`** package provides a single helper function, `GetEnv`, that reads an environment variable by key and logs a fatal error if the value is missing. It is intended to be used wherever configuration values need to be fetched from the process environment.

---

## Environment Variables & Configuration  

| Variable | Purpose |
|----------|---------|
| *Any* (`key string`) | The name of the environment variable that `GetEnv` will read. No default key is hard‑coded; callers supply it at runtime.|

No command‑line flags or arguments are defined in this file, but the function can be called from any other package (e.g., a CLI entry point) to obtain configuration values.

---

## Project Package Structure  

```
examples/
└── codemonkey/
    └── pkg/
        └── util/
            └── utils.go
```

* `utils.go` – contains the implementation of `GetEnv`.

---

## Code Relations & Edge Cases  

### Function `GetEnv`
- **Signature**: `func GetEnv(key string) string`
- **Behavior**  
  1. Calls `os.Getenv(key)` to fetch the value for the supplied key.  
  2. If the returned string is empty, it logs a fatal message via `log.Fatal().Msgf`.  
  3. Returns the fetched value.

### Edge Cases
- **Missing Key**: If the environment variable is not set, the function will log a fatal error and terminate the program.  
- **Empty Value**: The same fatal path applies if the key exists but holds an empty string.  

The package can be used directly from a command‑line entry point (e.g., `main.go`) by importing `"examples/codemonkey/pkg/util"` and calling `util.GetEnv("MY_VAR")`. No additional flags or arguments are required for this helper.

---