# Scenario

**Feature**: no target directory argument is supplied to fill-code

```
# agent reads requirement, invokes Fake Codex, writes output
doctest agent <cmd> --requirement req.md -> Fake Codex -> generated code

# session state tracked in event files
doctest <- Fake Codex (session id, events, progress)
```

## Preconditions
- No target directory argument is supplied to fill-code.

## Steps
1. Run `doctest agent fill-code`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "fill-code"}
    return nil
}
```

