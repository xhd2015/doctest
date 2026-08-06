# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- `--requirement` points to a nonexistent file, no prompt.

## Steps
1. Append `--requirement` with a nonexistent path (no file created).
2. StdinSource = "devnull" (terminal-like).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "devnull"
    req.Args = append(req.Args, "--requirement", "/nonexistent/path/requirement.md")
    return nil
}
```
