## Group: path-prefix
Tests for `./<prefix>/...` pattern support.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Logf("path-prefix group: WorkDir=%s", req.WorkDir)
    return nil
}
```
