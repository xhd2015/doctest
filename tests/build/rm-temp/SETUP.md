# Scenario

**Feature**: `doctest build` with `--rm` should delete the generated temp directory

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- `doctest build` with `--rm` should delete the generated temp directory.

## Steps
1. Run `doctest build <exampleDir> --rm`.
2. Parse the temp directory from stderr.
3. Verify the temp directory has been removed.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", exampleDir, "--rm"}
    return nil
}
```
