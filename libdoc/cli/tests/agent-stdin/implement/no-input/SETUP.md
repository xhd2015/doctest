## Preconditions
- No positional args, no stdin data.
- Stdin is set to `/dev/null` (character device, not a pipe).

## Steps
1. StdinSource = "devnull", no Stdin content, no positional args.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.StdinSource = "devnull"
    return nil
}
```
