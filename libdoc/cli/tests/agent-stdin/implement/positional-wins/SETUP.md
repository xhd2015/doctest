## Preconditions
- Both a positional arg and stdin data are present.

## Steps
1. StdinSource = "pipe", Stdin = "this should be ignored", Args = ["use positional instead"].

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "pipe"
    req.Stdin = "this stdin content should be ignored"
    req.Args = []string{"use positional instead"}
    return nil
}
```
