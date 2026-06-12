## Preconditions
- A temp project exists with `test_a`, `test_b`, and `test_c` doctest trees.

## Steps
1. Create the temp project.
2. Run `doctest test test_a test_b test_c` (three plain directory arguments).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    projDir := createMultiDirProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "test_a", "test_b", "test_c"}
    return nil
}
```
