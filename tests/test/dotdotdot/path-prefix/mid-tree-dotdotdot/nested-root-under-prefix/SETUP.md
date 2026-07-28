# Scenario

**Case**: nested DOCTEST under prefix found (`./tree/mid/...`)

Same full fixture as siblings-excluded (mid + nested + sibling).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createMidTreePrefixProject(t)
	req.Args = []string{"test", "-v", "--label-all", "./tree/mid/..."}
	return nil
}
```
