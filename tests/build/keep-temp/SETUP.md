## Preconditions
- `doctest build` without `--rm` should keep the generated temp directory.

## Steps
1. Run `doctest build <exampleDir>` (no --rm, no --gen-dir).
2. Parse the temp directory from stderr.
3. Verify the temp directory still exists.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", exampleDir}
    return nil
}
```
