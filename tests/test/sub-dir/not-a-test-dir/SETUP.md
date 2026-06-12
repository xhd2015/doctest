## Preconditions
- A plain directory is not part of any doc-style test tree (no DOCTEST.md, no SETUP.md anywhere).

## Steps
1. Create an empty plain directory.
2. Run `doctest test <plainDir>`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    plainDir := t.TempDir()
    req.Args = []string{"test", plainDir}
    return nil
}
```
