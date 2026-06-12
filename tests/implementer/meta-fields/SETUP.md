## Preconditions
- The implementer SETUP.md provides shared infrastructure.

## Steps
1. Test meta.json field structure and session persistence across calls.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=meta-fields")
    return nil
}
```
