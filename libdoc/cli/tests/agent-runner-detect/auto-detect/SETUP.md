## Preconditions
- This group tests branch: no `--agent-runner` flag → auto-detection kicks in.
- All children in this group set env vars that influence auto-detection.

## Steps
1. Prepend `"agent"`, `"implement"`, `"test"` to args (no `--agent-runner` flag).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = append([]string{"agent", "implement", "test"})
    return nil
}
```
