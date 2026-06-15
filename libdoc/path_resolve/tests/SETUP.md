## Preconditions
- The path_resolve package is importable from this test tree.
- The root Run dispatches to the function under test based on runType.

## Steps
1. Each leaf sets input fields and runType via `Setup`.
2. The root Run calls the corresponding function.
3. Each leaf asserts results via `Assert`.

## Context
- These are direct unit tests, not CLI integration tests.
- The root Run dispatches to the appropriate test function.

```go
import (
    "testing"

    "github.com/xhd2015/doctest/libdoc/path_resolve"
)

type Request struct {
    Input    string
    BasePath string
}

type Response struct {
    BoolResult   bool
    StringResult string
    DirsResult   []string
    RootResult   string
    RootOkResult bool
    ErrResult    string
}

var runType string

func Run(t *testing.T, req *Request) (*Response, error) {
    switch runType {
    case "extract_base_path":
        result := path_resolve.ExtractBasePath(req.Input)
        return &Response{StringResult: result}, nil
    case "find_dot_dot_dot_dirs":
        dirs, err := path_resolve.FindDotDotDotDirs(req.BasePath)
        if err != nil {
            return &Response{ErrResult: err.Error()}, err
        }
        return &Response{DirsResult: dirs}, nil
    case "is_dot_dot_dot_pattern":
        result := path_resolve.IsDotDotDotPattern(req.Input)
        return &Response{BoolResult: result}, nil
    case "resolve_root":
        root, ok := path_resolve.ResolveRoot(req.Input)
        return &Response{RootResult: root, RootOkResult: ok}, nil
    default:
        t.Fatalf("unknown runType: %s", runType)
        return nil, nil
    }
}
```
