## Preconditions
- `doctest build` with `--rm` should delete the generated temp directory.

## Steps
1. Run `doctest build <exampleDir> --rm`.
2. Parse the temp directory from stderr.
3. Verify the temp directory has been removed.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", exampleDir, "--rm"}
    return nil
}
```
