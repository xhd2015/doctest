## Preconditions
- A project with nested module at `sub/` whose module path `testproj2/sub` is NOT a child of parent `testproj`.

## Steps
1. Create temp project via cross-go-mod helper.
2. Run `doctest test -v ./...` from project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createCrossModuleProject(t, "sub", "testproj2/sub", "hidden_test")
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
