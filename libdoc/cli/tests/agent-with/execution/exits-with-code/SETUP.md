## Preconditions
- `sh` is available to exit with a specific code.

## Steps
1. Run a child that exits with code 42.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = append(req.Args, "sh", "-c", "exit 42")
    return nil
}
```
