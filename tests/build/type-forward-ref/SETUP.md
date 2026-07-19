# Scenario

**Feature**: a doc-style test tree where types are defined in forward-reference order

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree where types are defined in forward-reference order
  (LocationEntry references GitInfo before GitInfo is defined).
- The generator emits type declarations in file-order inside a function body,
  where Go requires dependencies to be defined first.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import (
"path/filepath"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"test", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "type-forward-ref")}
    return nil
}
```
