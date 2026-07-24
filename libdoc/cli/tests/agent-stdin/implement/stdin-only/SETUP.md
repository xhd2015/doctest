## Preconditions
- Stdin is a pipe with content, no positional args.

## Steps
1. StdinSource = "pipe", Stdin = "hello from stdin", no positional args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "pipe"
    req.Stdin = "hello from stdin"
    return nil
}
```
