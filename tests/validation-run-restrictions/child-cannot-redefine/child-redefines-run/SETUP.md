# Scenario

**Feature**: the `tests/testdata/child-redefines-run/` fixture has a root with Run

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- The `tests/testdata/child-redefines-run/` fixture has a root with Run
  and a leaf SETUP.md that also defines Run, which should be rejected.

## Steps
1. Run `doctest test` on the child-redefines-run fixture.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    fixtureDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "child-redefines-run")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
