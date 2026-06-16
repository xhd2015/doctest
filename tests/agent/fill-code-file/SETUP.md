# Scenario

**Feature**: the target path exists but is not a directory

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

## Preconditions
- The target path exists but is not a directory.

## Steps
1. Run `doctest agent fill-code <file>`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    path := filepath.Join(t.TempDir(), "target.txt")
    if err := os.WriteFile(path, []byte("not a dir"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"agent", "fill-code", path}
    return nil
}
```
