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
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    fixtureDir := filepath.Join(DOCTEST_ROOT, "testdata", "root-missing-run")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
