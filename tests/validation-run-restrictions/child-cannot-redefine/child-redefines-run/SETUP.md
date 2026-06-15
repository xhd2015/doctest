## Preconditions
- The `tests/testdata/child-redefines-run/` fixture has a root with Run
  and a leaf SETUP.md that also defines Run, which should be rejected.

## Steps
1. Run `doctest test` on the child-redefines-run fixture.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    fixtureDir := filepath.Join(DOCTEST_ROOT, "testdata", "child-redefines-run")
    req.Args = []string{"test", fixtureDir}
    return nil
}
```
