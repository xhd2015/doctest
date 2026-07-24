## Preconditions
- Simple `echo` command is available in PATH.

## Steps
1. Run `echo hello`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "echo", "hello")
    return nil
}
```
