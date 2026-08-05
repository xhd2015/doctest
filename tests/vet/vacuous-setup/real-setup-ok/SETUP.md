# Scenario

**Feature**: non-root SETUP whose `func Setup` does real work (sets a field) passes vet

```
# fixture leaf SETUP sets a field then returns
write DOCTEST + leaf/SETUP (req.InputDir = …)
  -> runner.VetArgs(["vet", dir])
  -> exit 0
```

## Preconditions

- Fixture tree under `t.TempDir()` with root `DOCTEST.md` and `leaf/SETUP.md`.
- Leaf Setup performs real work (field assign of a non-blank identifier target).

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
	leafDir := filepath.Join(dir, "leaf")
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), fixture, 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", dir}
	return nil
}
```
