## Preconditions
- A project with complex doctest tree structure (ancestor/ leaf/ nested sub2/).

## Steps
1. Create the mixed test project.
2. Run `doctest test -v ./ancestor/leaf` (without `...`) from the project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createMixedTestProject(t)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./ancestor/leaf"}
    return nil
}
```
