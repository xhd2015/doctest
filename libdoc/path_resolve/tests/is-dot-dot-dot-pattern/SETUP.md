# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- The `IsDotDotDotPattern` function is defined in `path_resolve`.
- Input is a string argument (e.g. `./foo/...`, `foo/bar`, `...`).

## Steps
1. Call `path_resolve.IsDotDotDotPattern(req.Input)`.
2. Return the result in `resp.BoolResult`.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    _ = req
    req.RunType = "is_dot_dot_dot_pattern"
    return nil
}
```
