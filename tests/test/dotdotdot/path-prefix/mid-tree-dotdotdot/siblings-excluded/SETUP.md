# Scenario

**Case**: no siblings outside prefix (`./tree/mid/...`)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createMidTreePrefixProject(t)
	// -count=1 disables leaf-cache skip so Setup MARKER logs always appear.
	req.Args = []string{"test", "-v", "--label-all", "-count=1", "./tree/mid/..."}
	return nil
}
```
