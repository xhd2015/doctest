## Preconditions
- The format string is `"\n"` (empty message that already ends with newline).

## Steps
1. Set `LOGF_FORMAT=\n` via env.
2. Call `Logf("\n")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=\n")
    return nil
}
```
