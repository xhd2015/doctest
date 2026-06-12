## Preconditions
- A doc-style test tree where two ancestor SETUP.md files define a helper
  function with the same name (e.g., `func myHelper` in both root and child).
- The generator emits both as closures using `:=`, causing "no new variables on
  left side of :=" because the second assignment reuses a variable name.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", filepath.Join(DOCTEST_ROOT, "build", "testdata", "helper-shadow")}
    return nil
}
```
