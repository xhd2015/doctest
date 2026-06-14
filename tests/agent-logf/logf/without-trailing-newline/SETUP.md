## Preconditions
- The format string `"hello"` has no trailing newline.

## Steps
1. Set `LOGF_FORMAT=hello` via env.
2. Call `Logf("hello")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=hello")
    return nil
}
```
