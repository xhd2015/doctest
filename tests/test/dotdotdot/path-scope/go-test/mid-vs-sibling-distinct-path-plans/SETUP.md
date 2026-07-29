# Scenario

**Case**: mid vs sibling scopes must produce **distinct, path-scoped** go test plans.

Not the same `go test ./__workspace/suite` for both — filter lives at go-test
level as path patterns under the selected subpath.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createPathScopeMidSibling(t)
	genDir := filepath.Join(req.WorkDir, ".gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Default Run is unused; Assert runs mid and sibling with shared gen-dir.
	req.Args = []string{"test", "-v", "--gen-dir", genDir, "-count=1", "./tree/mid/..."}
	return nil
}
```
