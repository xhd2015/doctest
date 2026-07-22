# Scenario

**Feature**: in-process discover + label filter selection (library contract)

```
# temp fixture mod (fast, slow, ui, both, heavy)
DiscoverTreeCasesLight(root)
  -> FilterBySubDir? (optional)
  -> FilterCasesByLabel(opts)
  -> run paths + skipped paths/reasons
```

## Preconditions

- Nested L2 root: no product binary; uses `libdoc/core` only.
- Fixture trees written under `t.TempDir()` per leaf.

## Steps

1. Leaf Setup writes the five-leaf fixture (or subset) and sets LabelExprs / SubDir.
2. Run discovers light + filters.
3. Assert compares sorted run/skipped paths (and reasons when relevant).

```go
import (
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

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}

func writeMinimalLeafSetupAssert(t *testing.T, root, leafRel string) {
	t.Helper()
	var setup strings.Builder
	setup.WriteString("# Scenario\n\n**Feature**: label-filter select fixture leaf\n\n")
	setup.WriteString(bt(3) + "\nfixture leaf\n" + bt(3) + "\n\n## Steps\n1. leaf setup\n\n")
	setup.WriteString(goFence())
	setup.WriteString("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
	setup.WriteString(endFence())
	var assert strings.Builder
	assert.WriteString("## Expected\n- passes\n\n")
	assert.WriteString(goFence())
	assert.WriteString("func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
	assert.WriteString(endFence())
	testtree.WriteFile(t, root, leafRel+"/SETUP.md", setup.String())
	testtree.WriteFile(t, root, leafRel+"/ASSERT.md", assert.String())
}

func writeLabeledAssert(t *testing.T, root, leafRel, label, explanation string) {
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
	b.WriteString("func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
	b.WriteString(endFence())
	testtree.WriteFile(t, root, leafRel+"/ASSERT.md", b.String())
}

// writeLabelFilterMod builds the standard five-leaf fixture (same shape as CLI).
func writeLabelFilterMod(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeafSetupAssert(t, root, "fast")
	writeMinimalLeafSetupAssert(t, root, "slow")
	writeLabeledAssert(t, root, "slow", "slow", "slow profile")
	writeMinimalLeafSetupAssert(t, root, "ui")
	writeLabeledAssert(t, root, "ui", "ui-automation", "browser ui")
	writeMinimalLeafSetupAssert(t, root, "both")
	writeLabeledAssert(t, root, "both", "slow, ui-automation", "slow ui combo")
	writeMinimalLeafSetupAssert(t, root, "heavy")
	writeLabeledAssert(t, root, "heavy", "heavy", "heavy profile")
	return root
}

func equalPaths(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func requirePaths(t *testing.T, got, want []string, label string) {
	t.Helper()
	if !equalPaths(got, want) {
		t.Fatalf("%s paths mismatch\ngot:  %v\nwant: %v", label, got, want)
	}
}

```
