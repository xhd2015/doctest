## Preconditions
- Tests in this group test `doctest agent implement --list-sessions`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-list")
    return nil
}
```
