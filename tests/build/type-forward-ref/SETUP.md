## Preconditions
- A doc-style test tree where types are defined in forward-reference order
  (LocationEntry references GitInfo before GitInfo is defined).
- The generator emits type declarations in file-order inside a function body,
  where Go requires dependencies to be defined first.

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", filepath.Join(DOCTEST_ROOT, "build", "testdata", "type-forward-ref")}
    return nil
}
```
