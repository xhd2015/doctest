# Scenario

**Feature**: `--timeout` is forwarded to the generated `go test` command

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- A valid doc-style test tree exists in the repository.

## Steps
1. Run `doctest test -v --timeout 45s <dir>`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"test", "-v", "--timeout", "45s", exampleDir}
    return nil
}
```