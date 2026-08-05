# Scenario

**Feature**: full-tree `doctest vet` enforces L3 (label `e2e`) share budget after structure + anti-patterns

```
# full vet only (not --changed)
DiscoverTreeCasesLight / inventory-aligned labels
  -> L3 = leaves with label "e2e"; L2 = all other leaves
  -> if leaves >= 10 && 100*L3/leaves > 10 → hard fail
  -> MaxL3Pct=10, MinLeaves=10 (fixed; no CLI flags this phase)
```

## Preconditions

- Fixtures are valid DOCTEST trees under `t.TempDir()` (DOCTEST + N leaf SETUP/ASSERT).
- L3 identity matches `doctest list`: **`label: e2e` only** (`heavy` alone is L2).
- In-process via shared root `Run` → `runner.VetArgs` (no product binary).
- Parallel-safe: no Setenv/Chdir/stdio reassignment; relative Args not required (abs dir).
- **`--changed` skips share**: no leaf here — real `--changed` needs git context and is
  covered elsewhere. Implementer **must** skip L3 share when `opts.ChangedOnly`
  (share is full-tree inventory only).

## Steps

1. Grouping Setup defines fixture helpers for labeled multi-leaf trees.
2. Each leaf writes a tailored tree and sets `req.Args = ["vet", dir]`.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Helpers only — no func Setup (organization + shared fixture writers).
// Leaves call writeShareFixture / shareSpecs / heavySpecs from this package scope.

// fixtureLeafSETUP is a vet-valid leaf SETUP.md: must start with "# Scenario"
// (WriteMinimalRunnableTree only emits "## Steps" and fails checkSETUPSections).
// Prose-only (no Go Setup): organization/fixture leaves need no harness Setup.
// Fence markers use hex escapes so this root SETUP does not close its own go fence.
const fixtureLeafSETUP = "# Scenario\n\n**Feature**: layer-share fixture leaf\n\n" +
	"\x60\x60\x60\n# fixture pipeline\nsystem -> run\n\x60\x60\x60\n\n" +
	"## Steps\n1. fixture leaf setup (prose-only; no Go Setup)\n"

// fixtureLeafASSERTBody is a vet-valid ASSERT.md body (Go Assert required by ParseAssertDocument).
const fixtureLeafASSERTBody = "## Expected\n- fixture leaf\n\n" +
	"\x60\x60\x60go\nimport (\n\t\"testing\"\n\n\t\"github.com/xhd2015/doctest/session\"\n)\n\n" +
	"func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {\n" +
	"\t_ = d\n}\n\x60\x60\x60\n"

// writeShareFixture writes a minimal valid DOCTEST tree with optional ASSERT labels.
// Each entry is "rel" or "rel|labelField" (e.g. "e2e0|e2e", "slow0|heavy").
// Leaf SETUP.md always starts with "# Scenario" so full vet structure checks pass.
func writeShareFixture(t *testing.T, root string, specs []string) {
	t.Helper()
	if len(specs) == 0 {
		t.Fatal("writeShareFixture: need at least one leaf spec")
	}
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.VetDOCTEST())
	for _, spec := range specs {
		rel, labels, _ := strings.Cut(spec, "|")
		if rel == "" {
			t.Fatalf("empty leaf name in spec %q", spec)
		}
		leafDir := filepath.Join(filepath.FromSlash(rel))
		testtree.WriteFile(t, root, filepath.Join(leafDir, "SETUP.md"), fixtureLeafSETUP)
		assertBody := fixtureLeafASSERTBody
		if labels != "" {
			assertBody = "---\nlabel: " + labels + "\n---\n\n" + assertBody
		}
		testtree.WriteFile(t, root, filepath.Join(leafDir, "ASSERT.md"), assertBody)
	}
}

// shareSpecs builds total leaf specs: first e2eCount labeled "e2e", rest unlabeled.
func shareSpecs(total, e2eCount int) []string {
	if e2eCount > total {
		e2eCount = total
	}
	out := make([]string, 0, total)
	for i := 0; i < total; i++ {
		name := "leaf_" + itoaShare(i)
		if i < e2eCount {
			out = append(out, name+"|e2e")
		} else {
			out = append(out, name)
		}
	}
	return out
}

// heavySpecs builds total leaves all labeled "heavy" (no e2e → all L2).
func heavySpecs(total int) []string {
	out := make([]string, 0, total)
	for i := 0; i < total; i++ {
		out = append(out, "leaf_"+itoaShare(i)+"|heavy")
	}
	return out
}

func itoaShare(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
```
