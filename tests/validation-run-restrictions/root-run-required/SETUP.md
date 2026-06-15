## Preconditions
- Rule R1: Root SETUP.md must define type Request, type Response, and func Run.
- Without Run at the root, tree discovery fails.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    _ = req
    return nil
}
```
