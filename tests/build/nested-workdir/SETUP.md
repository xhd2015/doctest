# Scenario

**Feature**: a valid doc-style test tree exists in the module

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A valid doc-style test tree exists in the module.
- The command is invoked from inside a nested module directory.

## Steps
1. Set the process working directory to the module root (`DOCTEST_ROOT/..`).
2. Run `doctest build <absolute-dir>`.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata/basic-request-runner")
    req.WorkDir = filepath.Join(d.DOCTEST_ROOT, "..")
    req.Args = []string{"build", exampleDir}
    return nil
}
```
