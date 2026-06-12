## Preconditions
- A doc-style test tree where the root SETUP.md defines multiple helper functions,
  but a leaf test only uses a subset of them.
- The generator emits all helpers from every ancestor as closures.
  Unused closures cause "declared and not used" compilation errors.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", filepath.Join(DOCTEST_ROOT, "build", "testdata", "unused-helper")}
    return nil
}
```
