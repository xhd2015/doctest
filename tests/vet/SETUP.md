## Preconditions
- The vet command checks doc-style tree structure (renamed from validate).

## Steps
1. Choose a target directory.
2. Run `doctest vet`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    req.Env = append(req.Env, "DOCTEST_VET_TEST=1")
    return nil
}
```
