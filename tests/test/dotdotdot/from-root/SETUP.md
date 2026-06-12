## Preconditions
- A temporary project with go.mod + multiple DOCTEST.md trees exists.

## Steps
1. Create temp project via parent helper.
2. Run `doctest test ./...` from the project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createTempProject(t, req)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
