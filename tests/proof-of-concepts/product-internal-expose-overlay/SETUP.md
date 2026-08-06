# Scenario

**Feature**: materialize multi-module workspace + product-expose overlay

```
testdata/workspace -> t.TempDir()/ws
  app/ + runner/ + expose-src/ + overlay.json
leaf sets Scenario; WorkDir = runner
```

## Preconditions

- `go` on PATH.
- Fixture at `DOCTEST_ROOT/testdata/workspace`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if d == nil || d.DOCTEST_ROOT == "" {
		t.Fatal("DOCTEST_ROOT required")
	}
	fixture := filepath.Join(d.DOCTEST_ROOT, "testdata", "workspace")
	dest := filepath.Join(t.TempDir(), "workspace")
	materializeWorkspace(t, fixture, dest)
	req.WorkDir = filepath.Join(dest, "runner")
	req.AppRoot = filepath.Join(dest, "app")
	req.OverlayPath = filepath.Join(dest, "overlay.json")
	req.Scenario = ""
	return nil
}
```
