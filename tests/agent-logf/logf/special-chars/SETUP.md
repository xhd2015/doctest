## Preconditions
- The format string contains a line break, special characters, and already ends with `\n`.

## Steps
1. Set `LOGF_FORMAT=line1\nline2 -- special: !@#$%\n` via env.
2. Call `Logf("line1\nline2 -- special: !@#$%\n")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=line1\nline2 -- special: !@#$%\n")
    return nil
}
```
