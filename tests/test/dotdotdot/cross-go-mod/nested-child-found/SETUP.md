## Preconditions
- A project with nested module `testproj/sub` whose module path IS a child of parent `testproj`.

## Steps
1. Create temp project via cross-go-mod helper.
2. Run `doctest test -v ./...` from project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createCrossModuleProject(t, "sub", "testproj/sub", "child_test")
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
