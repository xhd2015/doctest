## Preconditions
- Unknown runner flags should fail.

## Steps
1. Run `doctest test <dir> --definitely-not-real`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"test", exampleDir, "--definitely-not-real"}
    return nil
}
```
