## Preconditions
- Tests in this group verify that `--trace` numbering is continuous (no gaps caused by non-displayable events).
- The fix moves `n++` inside the `if formatted != ""` block so only visible output consumes a number.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=agent-trace-numbering")
    return nil
}
```
