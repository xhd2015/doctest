# Scenario

**Feature**: the `tests/testdata/root-run-used/` fixture has a root SETUP.md with

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- The `tests/testdata/root-run-used/` fixture has a root SETUP.md with
  Request, Response, Setup, and Run — satisfying R1.
- A leaf with Setup and Assert completes the tree.

## Steps
1. Run `doctest test` on the valid root-run-used fixture.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    fixtureDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "root-run-used")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
