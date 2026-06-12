## Preconditions
- A valid doc-style test tree exists in the repository.

## Steps
1. Run `doctest test -v <dir>`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"test", "-v", exampleDir}
    return nil
}
```
