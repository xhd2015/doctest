# Scenario

**Feature**: `testdata/` is not an extra root and ASSERT under testdata is not counted

```
Harness -> root with own leaf + testdata/fake-root(DOCTEST) + testdata/ASSERT
  -> list <root>
  -> only the real root; leaves=1
```

## Preconditions

- Root has one real leaf.
- Under `testdata/`: a nested-looking DOCTEST.md and an ASSERT.md that must be ignored.

## Steps

1. Write root + real leaf.
2. Plant DOCTEST.md and ASSERT.md under testdata/.
3. Args = `list <root>`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "tree")
	writeLabeledLeaves(t, root, []string{"real"})
	// decoy under testdata: must not become a root or add leaves
	td := filepath.Join(root, "testdata", "decoy")
	if err := os.MkdirAll(td, 0755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteFile(t, td, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	testtree.WriteFile(t, td, "ASSERT.md", "## Expected\n- decoy\n")
	// also a loose ASSERT under testdata/
	testtree.WriteFile(t, filepath.Join(root, "testdata"), "ASSERT.md", "## Expected\n- ignored\n")
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
