## Preconditions
- Rule R2: Every non-root SETUP.md must have func Setup.
- Run is reserved for root; non-root SETUP.md without at least Setup is invalid.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    _ = req
    return nil
}
```
