## Preconditions
- Positional arg is present as the prompt.

## Steps
1. StdinSource = "devnull", Args = ["add dark mode support"].

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.StdinSource = "devnull"
    req.Args = []string{"add dark mode support"}
    return nil
}
```
