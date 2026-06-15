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
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    fixtureDir := filepath.Join(DOCTEST_ROOT, "testdata", "child-no-setup")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
