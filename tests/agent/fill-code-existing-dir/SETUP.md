## Preconditions
- The target directory exists.

## Steps
1. Run `doctest agent fill-code <dir>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "fill-code", t.TempDir()}
    return nil
}
```
