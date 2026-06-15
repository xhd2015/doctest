## Preconditions
- Rule R4: Only the root's Run is used in generated code.
- The root Run function is the one executed, not any overrides that may exist.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    _ = req
    return nil
}
```
