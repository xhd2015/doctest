## Preconditions
- No previous command arguments are required.

## Steps
1. Run `doctest --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"--help"}
    return nil
}
```

