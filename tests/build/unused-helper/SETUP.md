# Scenario

**Feature**: a doc-style test tree where the root SETUP.md defines multiple helper functions,

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree where the root SETUP.md defines multiple helper functions,
  but a leaf test only uses a subset of them.
- The generator emits all helpers from every ancestor as closures.
  Unused closures cause "declared and not used" compilation errors.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import (
"path/filepath"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"test", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "unused-helper")}
    return nil
}
```
