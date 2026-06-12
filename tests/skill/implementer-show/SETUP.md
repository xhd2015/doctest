## Preconditions
- The implementer prompt specification exists under `agents/doctest/libdoc/implementer`.

## Steps
1. Run `doctest skill implementer show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "implementer", "show"}
    return nil
}
```
