## Preconditions
- Rule R3: Non-root SETUP.md files cannot redefine func Run.
- Run is reserved for the root SETUP.md only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    _ = req
    return nil
}
```
