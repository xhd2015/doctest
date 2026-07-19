# Scenario

**Feature**: a valid doc-style test tree exists in the repository

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A valid doc-style test tree exists in the repository.

## Steps
1. Run `doctest build -v <dir>`.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
    req.Args = []string{"build", "-v", exampleDir}
    return nil
}
```
