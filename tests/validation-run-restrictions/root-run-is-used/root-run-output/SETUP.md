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
  Request{Name}, Response{Message}, and Run that returns a greeting.

## Steps
1. Run `doctest test` on the root-run-used fixture.
2. The root Run generates the greeting; the leaf Assert checks it.

## Context
- This verifies that only the root's Run is used in generated code.
  The leaf's Response must reflect the root Run's output.

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
