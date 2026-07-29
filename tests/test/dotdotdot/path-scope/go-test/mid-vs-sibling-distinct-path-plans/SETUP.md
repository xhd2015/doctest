# Scenario

**Case**: mid vs sibling scopes must produce **distinct, path-scoped** go test plans.

Expect `go test ./tree/mid/...` vs `./tree/sibling/...` — not shared
`./__workspace/suite` and not a hard-coded `*/suite` package.

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
