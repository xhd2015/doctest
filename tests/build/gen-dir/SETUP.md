## Preconditions
- A valid doc-style test tree exists in the repository.

## Steps
1. Run `doctest build <dir> --gen-dir <tmp>`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    genDir := filepath.Join(t.TempDir(), "generated")
    req.Args = []string{"build", exampleDir, "--gen-dir", genDir}
    return nil
}
```
