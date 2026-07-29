# Scenario

**Case**: `./tree/mid/...` go test runs mid leaf only (no sibling marker).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createPathScopeMidSibling(t)
	req.Args = []string{"test", "-v", "--label-all", "-count=1", "./tree/mid/..."}
	return nil
}
```
