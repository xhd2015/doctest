# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- The `ExtractBasePath` function is defined in `path_resolve`.
- Input is a pattern string like `"./foo/..."`.

## Steps
1. Call `path_resolve.ExtractBasePath(req.Input)`.
2. Return the result in `resp.StringResult`.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    _ = req
    req.RunType = "extract_base_path"
    return nil
}
```
