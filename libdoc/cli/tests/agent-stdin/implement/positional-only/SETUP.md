## Preconditions
- Positional arg is present as the prompt.

## Steps
1. Args include positional prompt text, stdin is terminal (devnull).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "devnull"
    req.Args = []string{"fix login bug"}
    return nil
}
```
