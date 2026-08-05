# Scenario

**Feature**: non-root SETUP whose body is only `_ = …` assigns then `return nil` fails with the same remove-code-block class

```
# fixture leaf SETUP is blank-assign + return nil
write DOCTEST + leaf/SETUP (_ = t; _ = d; _ = req; return nil)
  -> runner.VetArgs(["vet", dir])
  -> non-zero; same message class as pure return nil
  -/-> "implement the behavior"
```

## Preconditions

- Fixture tree under `t.TempDir()` with root `DOCTEST.md` and `leaf/SETUP.md`.
- Leaf Setup body is only blank-identifier assigns then `return nil`.
- Same message class as pure `return nil` (remove + code block).

## Steps

1. Write fixture under `t.TempDir()`.
2. Run `vet <dir>` in-process.

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
