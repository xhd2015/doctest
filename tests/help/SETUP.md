## Preconditions
- The doctest command supports top-level and scoped help output.

## Steps
1. Choose a help command variant.
2. Run the doctest command.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    req.Env = append(req.Env, "DOCTEST_HELP_TEST=1")
    return nil
}
```

