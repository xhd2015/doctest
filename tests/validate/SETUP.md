## Preconditions
- The validate command checks doc-style tree structure.

## Steps
1. Choose a target directory.
2. Run `doctest validate`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    req.Env = append(req.Env, "DOCTEST_VALIDATE_TEST=1")
    return nil
}
```

