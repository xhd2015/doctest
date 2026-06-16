# Scenario

**Feature**: the `tests/testdata/child-setup-only-valid/` fixture has a root with Run

```
# enforce SETUP.md rules at build time
root missing Run -> build error
child redefines Run -> build error
non-root missing Setup -> build error
```

## Preconditions
- The `tests/testdata/child-setup-only-valid/` fixture has a root with Run
  and a child with Setup only — the valid pattern for non-root SETUP.md.

## Steps
1. Run `doctest test` on the child-setup-only-valid fixture.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    fixtureDir := filepath.Join(DOCTEST_ROOT, "testdata", "child-setup-only-valid")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
