# Scenario

**Case**: cold `./tree/mid/...` generate must not emit sibling packages under gen.

```
gen-dir empty
doctest test --gen-dir <work>/.gen ./tree/mid/...   (bypass go test)
  -> .gen has tree/mid/...
  -> .gen must not have tree/sibling/...
```

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
	req.Env = []string{
		"GOWORK=off",
		"DOCTEST_CACHE_HOME=" + t.TempDir(),
		"DOCTEST_DEBUG=bypass-go-test=1",
	}
	req.Args = []string{"test", "-v", "--gen-dir", genDir, "-count=1", "./tree/mid/..."}
	return nil
}
```
