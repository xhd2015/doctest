# Scenario

**Feature**: the target directory exists

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

## Preconditions
- The target directory exists.

## Steps
1. Run `doctest agent fill-code <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "fill-code", t.TempDir()}
    return nil
}
```
