# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- `--requirement` flag points to a file with content, no positional or stdin prompt.

## Steps
1. Write a temp requirement file with content, append `--requirement` flag.
2. StdinSource = "devnull" (terminal-like), no positional args.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    tf := filepath.Join(t.TempDir(), "requirement.md")
    if err := os.WriteFile(tf, []byte("Design a login page with email and password fields"), 0644); err != nil {
        return err
    }
    req.Args = append(req.Args, "--requirement", tf)
    req.StdinSource = "devnull"
    return nil
}
```
