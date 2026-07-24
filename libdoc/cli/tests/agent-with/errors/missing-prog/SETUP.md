## Preconditions
- `--agent-runner` is provided with a value, but no `<prog>` argument.

## Steps
1. Args include `--agent-runner=opencode` but no prog.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "--agent-runner=opencode")
    return nil
}
```
