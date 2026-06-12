## Preconditions
- The implementer SETUP.md provides shared infrastructure (binaries and helpers).

## Steps
1. Test session ID resolution from the three sources in priority order.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=session-id")
    return nil
}
```
