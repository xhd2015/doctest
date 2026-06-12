## Preconditions
- Unknown runner flags should fail.

## Steps
1. Run `doctest build <dir> --definitely-not-real`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", exampleDir, "--definitely-not-real"}
    return nil
}
```
