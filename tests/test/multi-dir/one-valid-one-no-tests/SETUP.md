## Preconditions
- A temp project exists with one valid doctest tree (`test_a`) and one non-doctest directory (`no_tests`).

## Steps
1. Create the temp project.
2. Run `doctest test no_tests test_a` — `no_tests` has no tests and should be silently skipped, `test_a` should run and pass.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    projDir := createMultiDirProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "no_tests", "test_a"}
    return nil
}
```
