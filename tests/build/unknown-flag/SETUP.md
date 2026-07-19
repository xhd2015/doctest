# Scenario

**Feature**: unknown runner flags should fail

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- Unknown runner flags should fail.

## Steps
1. Run `doctest build <dir> --definitely-not-real`.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", exampleDir, "--definitely-not-real"}
    return nil
}
```
