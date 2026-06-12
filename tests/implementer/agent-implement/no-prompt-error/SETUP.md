## Preconditions
- No prompt argument is provided.

## Steps
1. Run `doctest agent implement` with no prompt.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement"}
    return nil
}
```
