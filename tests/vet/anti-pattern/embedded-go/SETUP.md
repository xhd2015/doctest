# Scenario

**Feature**: a doctest tree with a SETUP.md that embeds a Go program as a raw string literal

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A doctest tree with a SETUP.md that embeds a Go program as a raw string literal.

## Steps
1. Create a minimal doctest tree with a SETUP.md containing a string literal that has `package main` and `func main()`.
2. Run `doctest vet <dir>`.

```go
import (
"github.com/xhd2015/doctest/session"
    "github.com/xhd2015/doctest/libdoc/testtree"
	"os"
	"path/filepath"
	"testing"
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
