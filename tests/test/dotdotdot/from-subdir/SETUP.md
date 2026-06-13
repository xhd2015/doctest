## Preconditions
- The parent `dotdotdot` SETUP provides helpers (`createTestTree`, `createTempProject`).
- Each leaf creates its own temp project and sets WorkDir + Args.

## Steps
1. Create a temp project with appropriate DOCTEST.md trees.
2. Run `doctest test ./...` (or `./subpath/...`) from a subdirectory.
3. Verify discovery is scoped to the working directory, not the module root.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```