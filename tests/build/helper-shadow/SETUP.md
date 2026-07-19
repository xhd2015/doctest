# Scenario

**Feature**: a doc-style test tree where two ancestor SETUP.md files define a helper

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree where two ancestor SETUP.md files define a helper
  function with the same name (e.g., `func myHelper` in both root and child).
- The generator emits both as closures using `:=`, causing "no new variables on
  left side of :=" because the second assignment reuses a variable name.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import (
"path/filepath"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"test", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "helper-shadow")}
    return nil
}
```
