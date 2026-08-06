# Scenario

**Feature**: the `tests/testdata/root-missing-run/` fixture has a root SETUP.md with

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root may omit Setup (prose-only OK)
```

## Preconditions
- The `tests/testdata/root-missing-run/` fixture has a root SETUP.md with
  Request, Response, and Setup but deliberately no func Run.

## Steps
1. Run `doctest test` on the root-missing-run fixture.

## Context
- The doctest tool should report a validation error because the root
  SETUP.md must have func Run.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    fixtureDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "root-missing-run")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
