# Scenario

**Feature**: a doc-style test tree where a child SETUP.md's `func Setup` calls `Run(t, d, req)`

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree where a child SETUP.md's `func Setup` calls `Run(t, d, req)`.
- The generator lowers `func Run` to lowercase closure `run`, then aliases
  `Run := run` so Setup can call uppercase `Run` with the inject `d`.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import (
"path/filepath"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"test", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "call-run-from-setup")}
    return nil
}
```
