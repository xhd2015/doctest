# Scenario

**Feature**: the `tests/testdata/child-no-setup/` fixture has a leaf SETUP.md with

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root missing Setup -> build error
```

## Preconditions
- The `tests/testdata/child-no-setup/` fixture has a leaf SETUP.md with
  only type declarations and no func Setup or func Run.

## Steps
1. Run `doctest test` on the child-no-setup fixture.

## Context
- Non-root SETUP.md must have func Setup. Having neither Setup nor Run
  should trigger a validation error.

```go
import (
"github.com/xhd2015/doctest/session"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    fixtureDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "child-no-setup")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
