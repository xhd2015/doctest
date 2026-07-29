# Scenario

**Case**: nested go.mod under prefix with DOCTEST (policy B/C — include)

```
tree/mid/nestedmod/go.mod   module midtreeproj/nestedmod
tree/mid/nestedmod/suite/   DOCTEST → NESTED_MOD_LEAF
./tree/mid/...              must find nested-mod suite; no sibling
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createMidTreeNestedGomodProject(t)
	req.Args = []string{"test", "-v", "--label-all", "-count=1", "./tree/mid/..."}
	return nil
}
```
