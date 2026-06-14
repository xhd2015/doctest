## Preconditions
- The format string `"hello\n"` already ends with a newline.

## Steps
1. Set `LOGF_FORMAT=hello\n` via env.
2. Call `Logf("hello\n")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=hello\n")
    return nil
}
```
