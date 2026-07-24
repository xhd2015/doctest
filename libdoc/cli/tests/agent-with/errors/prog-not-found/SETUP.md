## Preconditions
- The specified `<prog>` does not exist in PATH.

## Steps
1. Args include `--agent-runner=opencode` and a nonexistent program name.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "--agent-runner=opencode", "nonexistent-program-12345")
    return nil
}
```
