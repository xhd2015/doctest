# Scenario

**Feature**: intermediate or leaf SETUP.md with `# Scenario` prose and no Go code block is accepted by vet

```
# prose-only SETUP (no go fence)
write DOCTEST + group/SETUP (Scenario prose only)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Fixture tree under `t.TempDir()` with root `DOCTEST.md` and intermediate
  `group/SETUP.md` that has `# Scenario` but **no** go code block.
- Intermediate SETUP without Go is allowed (same class as leaf prose-only).

## Steps

1. Write fixture under `t.TempDir()`.
2. Run `vet <dir>` in-process — expect exit 0.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.VetDOCTEST()), 0644); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile(filepath.Join(d.DOCTEST_CASE, "fixture_setup.md.txt"))
	if err != nil {
		t.Fatal(err)
	}
	groupDir := filepath.Join(dir, "group")
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(groupDir, "SETUP.md"), fixture, 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", dir}
	return nil
}
```
