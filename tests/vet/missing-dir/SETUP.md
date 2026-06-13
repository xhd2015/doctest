## Preconditions
- No target directory argument is supplied.

## Steps
1. Run `doctest vet`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"vet"}
    return nil
}
```
