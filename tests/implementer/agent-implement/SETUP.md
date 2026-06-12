## Preconditions
- Tests in this group test `doctest agent implement`.
- The doctest binary and fake-codex are already built by parent.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-implement")
    return nil
}
```
