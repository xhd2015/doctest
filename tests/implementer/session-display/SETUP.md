## Preconditions
- The implementer SETUP.md provides shared infrastructure (binaries and helpers).

## Steps
1. Test session display messages for created and resumed sessions.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=session-display")
    return nil
}
```
