# Scenario

**Case**: full-tree gen, then mid-only gen — content outside `tree/mid` must stay
byte-identical (except bookkeeping).

Both CLI phases run inside Assert (shared Request has one argv). Locks generate
phase must not rewrite sibling / suite / workspace outside the mid path.

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
	// Dummy Args so Run is a no-op soft path; Assert drives both gen phases.
	req.Env = []string{
		"GOWORK=off",
		"DOCTEST_CACHE_HOME=" + t.TempDir(),
		"DOCTEST_DEBUG=bypass-go-test=1",
	}
	req.Args = []string{"test", "-v", "--gen-dir", genDir, "-count=1", "./tree/mid/..."}
	return nil
}
```
