## Preconditions
- The `ResolveRoot` function is defined in `path_resolve`.
- It walks up the directory tree to find the root containing DOCTEST.md (or falls back to SETUP.md).

## Steps
1. Call `path_resolve.ResolveRoot(req.Input)`.
2. Return the result in `resp.RootResult` and `resp.RootOkResult`.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    _ = req
    req.RunType = "resolve_root"
    return nil
}
```
