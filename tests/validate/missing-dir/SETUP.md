## Preconditions
- No target directory argument is supplied.

## Steps
1. Run `doctest validate`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"validate"}
    return nil
}
```

