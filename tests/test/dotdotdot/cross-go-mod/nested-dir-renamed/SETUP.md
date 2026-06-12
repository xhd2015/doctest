## Preconditions
- A project with nested module at `some/` whose module path `testproj/some/renamed` IS a child of parent `testproj`.
- The directory name (`some`) differs from the module path suffix (`some/renamed`).

## Steps
1. Create temp project via cross-go-mod helper.
2. Run `doctest test -v ./...` from project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createCrossModuleProject(t, "some", "testproj/some/renamed", "child_test")
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
