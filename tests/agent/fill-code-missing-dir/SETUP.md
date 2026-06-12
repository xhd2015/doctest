## Preconditions
- No target directory argument is supplied to fill-code.

## Steps
1. Run `doctest agent fill-code`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "fill-code"}
    return nil
}
```

