# Scenario

**Feature**: a doctest tree with a SETUP.md that only reads/writes via `os.Stdout` without reassignment (positive control)

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# allowed
fmt.Fprint(os.Stdout, ...) | read os.Stdout/os.Stderr/os.Stdin without assign
```

## Preconditions
- A doctest tree with a SETUP.md Go block that uses `fmt.Fprint(os.Stdout, ...)` only (no `os.Stdout =`).

## Steps
1. Create a minimal doctest tree with a SETUP.md that writes to `os.Stdout` without reassignment.
2. Run `doctest vet <dir>` in-process — expect exit 0.

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
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), fixture, 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", dir}
	return nil
}
```
