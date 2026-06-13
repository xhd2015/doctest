## Preconditions
- Tests in this group test `doctest agent implement --trace`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-trace")
    return nil
}
```
