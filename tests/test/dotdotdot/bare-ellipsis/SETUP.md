## Preconditions
- The parent dotdotdot helper functions are available.

## Steps
1. Run `doctest test ...` (bare ellipsis, no `./` prefix, no path).
2. This should produce an error because bare `...` is not supported.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = t.TempDir()
    req.Args = []string{"test", "-v", "..."}
    return nil
}
```
