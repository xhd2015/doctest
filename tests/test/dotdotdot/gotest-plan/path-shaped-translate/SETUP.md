# Scenario

**Feature**: path-shaped TranslatePath contract (mid + nested go.mod)

Does not require path-shaped **execution** yet: locks pure gotestmap rules used
when ModePathShaped is selected. CLI smoke uses a tiny no-op path.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Minimal CLI touch so leaf is a full doctest run; plan matrix is pure Assert.
	req.WorkDir = createSingleModTwoTrees(t)
	req.Args = []string{"test", "-v", "./alpha/..."}
	return nil
}
```
