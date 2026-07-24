# Scenario

**Feature**: in-process discovery skip / label-all / frontmatter policy

```
# temp fixture tree
DiscoverTreeCasesLight(root)
  -> FilterBySubDir? 
  -> FilterCasesByLabel / PartitionLabeledCases
  -> run + skipped paths

# pure frontmatter
ParseAssertFrontmatter(content) -> labels, explanation | error
```

## Preconditions

- Nested L2 root: `libdoc/core` only; no product binary.
- Fixtures under `t.TempDir()`.

## Steps

1. Leaf Setup builds fixture or frontmatter string and sets Op/options.
2. Run discovers+filters or parses.
3. Assert path sets or parse fields.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func bt(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '`'
	}
	return string(b)
}

func goFence() string {
	return bt(3) + "go\n"
}

func endFence() string {
	return bt(3) + "\n"
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}

func writeMinimalLeaf(t *testing.T, root, leafRel string) {
	t.Helper()
	var setup strings.Builder
	setup.WriteString("# Scenario\n\n**Feature**: label-skip select fixture leaf\n\n")
	setup.WriteString(bt(3) + "\nfixture leaf\n" + bt(3) + "\n\n## Steps\n1. leaf setup\n\n")
	setup.WriteString(goFence())
	setup.WriteString("import \"testing\"\n\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n")
	setup.WriteString(endFence())
	var assert strings.Builder
	assert.WriteString("## Expected\n- passes\n\n")
	assert.WriteString(goFence())
	assert.WriteString("func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n")
	assert.WriteString(endFence())
	testtree.WriteFile(t, root, leafRel+"/SETUP.md", setup.String())
	testtree.WriteFile(t, root, leafRel+"/ASSERT.md", assert.String())
}

func writeLabeledAssert(t *testing.T, root, leafName, label, explanation string) {
	t.Helper()
	var b strings.Builder
	b.WriteString("---\n")
	if label != "" {
		b.WriteString("label: ")
		b.WriteString(label)
		b.WriteString("\n")
	}
	if explanation != "" {
		b.WriteString("explanation: ")
		b.WriteString(explanation)
		b.WriteString("\n")
	}
	b.WriteString("---\n\n## Expected\n- passes\n\n")
	b.WriteString(goFence())
	b.WriteString("func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n")
	b.WriteString(endFence())
	assertPath := filepath.Join(root, leafName, "ASSERT.md")
	if err := os.WriteFile(assertPath, []byte(b.String()), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeLabeledTree(t *testing.T, includeFast bool, label, explanation string) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	if includeFast {
		writeMinimalLeaf(t, root, "fast_leaf")
	}
	writeMinimalLeaf(t, root, "labeled_leaf")
	writeLabeledAssert(t, root, "labeled_leaf", label, explanation)
	return root
}

func writeExplanationOnlyTree(t *testing.T, explanation string) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeaf(t, root, "explained_leaf")
	writeLabeledAssert(t, root, "explained_leaf", "", explanation)
	return root
}

func writeUnlabeledTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeaf(t, root, "plain_leaf")
	return root
}

func writeGroupingLabeledTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	var groupSetup strings.Builder
	groupSetup.WriteString("# Scenario\n\n**Feature**: e2e grouping fixture\n\n")
	groupSetup.WriteString(bt(3) + "\ne2e grouping\n" + bt(3) + "\n\n## Steps\n1. e2e grouping node\n\n")
	groupSetup.WriteString(goFence())
	groupSetup.WriteString("import \"testing\"\n\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n")
	groupSetup.WriteString(endFence())
	testtree.WriteFile(t, root, "e2e/SETUP.md", groupSetup.String())
	writeMinimalLeaf(t, root, "e2e/fast_child")
	writeMinimalLeaf(t, root, "e2e/labeled_child")
	writeLabeledAssert(t, root, "e2e/labeled_child", "ui-automation", "grouping skip")
	return root
}

func requirePaths(t *testing.T, got, want []string, label string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s paths mismatch\ngot:  %v\nwant: %v", label, got, want)
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("%s paths mismatch\ngot:  %v\nwant: %v", label, got, want)
			return
		}
	}
}
```
