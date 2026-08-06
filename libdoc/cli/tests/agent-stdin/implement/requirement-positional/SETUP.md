# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- `--requirement` flag with file content and a positional arg.

## Steps
1. Write a temp requirement file, append `--requirement` flag + positional arg.
2. StdinSource = "devnull" (terminal-like).

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    tf := filepath.Join(t.TempDir(), "requirement.md")
    if err := os.WriteFile(tf, []byte("Feature spec: toggle theme"), 0644); err != nil {
        return err
    }
    req.Args = append(req.Args, "--requirement", tf)
    req.Args = append(req.Args, "add dark mode toggle")
    req.StdinSource = "devnull"
    return nil
}
```
