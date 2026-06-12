## Preconditions
- A temp project exists with only a non-doctest directory (`no_tests`). No valid test trees are targeted.

## Steps
1. Create the temp project.
2. Run `doctest test no_tests` — the dir has no doctest tree, so `ErrNoTestsFound` is returned.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    projDir := createMultiDirProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "no_tests"}
    return nil
}
```
