## Preconditions
- This group tests branch: user explicitly passes `--agent-runner` flag → no auto-detection, value passed through directly.

## Steps
1. Prepend `"agent"`, `"implement"`, `"test"`, and `"--agent-runner=idonotexist"` to args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append([]string{"agent", "implement", "test", "--agent-runner=idonotexist"})
    return nil
}
```
