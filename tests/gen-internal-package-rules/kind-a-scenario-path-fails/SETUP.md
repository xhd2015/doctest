# Scenario

**Feature**: Kind A — scenario path `http/internal/…` breaks unified gen suite packaging

```
doctest test <fixture with leaf under http/internal/post-succeeds>
  -> mapping-gen / __allleaves imports testcase/…/http/internal/…
  -> FAIL use of internal package
```

## Steps

1. Copy `testdata/scenario-path-internal` into a temp project module.
2. Run product `doctest test` on that tree (default external unified gen).

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	fixture := mustCopyFixture(t, d, "scenario-path-internal")
	// Parent module is unrelated; leaf Run does not import product internal.
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module example.com/kind-a-proj\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	tree := filepath.Join(proj, "tests")
	if err := copyDir(fixture, tree); err != nil {
		t.Fatalf("place tree: %v", err)
	}
	req.Kind = "A"
	req.WorkDir = proj
	req.Args = []string{"test", "-count=1", "-v", "./tests/..."}
	return nil
}
```
