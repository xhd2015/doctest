# Scenario

**Feature**: labeled ASSERT.md frontmatter controls discovery-mode skip behavior

```
# build doctest binary once per run
Setup -> build ./cmd/doctest -> req.Bin

# subprocess exercises CLI contract
req.Bin <args> -> stdout/stderr/exit captured in Response
```

## Preconditions

- Fresh doctest binary built from module source in this root Setup.

## Steps

1. Build a temp doctest tree matching the leaf scenario.
2. Configure `req.Args` and optional `req.WorkDir`.
3. Assert complete observable output blocks.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m")
var finalSummaryPlainRe = regexp.MustCompile("^(PASS|FAIL) \\(\\d+/\\d+\\) in .+$")

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
	req.Timeout = 120 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", "..", ".."))
	return nil
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
	b.WriteString("func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
	b.WriteString(endFence())
	assertPath := filepath.Join(root, leafName, "ASSERT.md")
	if err := os.WriteFile(assertPath, []byte(b.String()), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeMinimalLeaf(t *testing.T, root, leafRel string) {
	t.Helper()
	var setup strings.Builder
	setup.WriteString("# Scenario\n\n**Feature**: fixture leaf for label-skip tests\n\n")
	setup.WriteString(bt(3) + "\nfixture leaf\n" + bt(3) + "\n\n## Steps\n1. leaf setup\n\n")
	setup.WriteString(goFence())
	setup.WriteString("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
	setup.WriteString(endFence())
	setup.WriteString("\n")
	var assert strings.Builder
	assert.WriteString("## Expected\n- passes\n\n")
	assert.WriteString(goFence())
	assert.WriteString("func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
	assert.WriteString(endFence())
	assert.WriteString("\n")
	testtree.WriteFile(t, root, leafRel+"/SETUP.md", setup.String())
	testtree.WriteFile(t, root, leafRel+"/ASSERT.md", assert.String())
}

func writeLabeledTree(t *testing.T, includeFast bool, label string, explanation string) string {
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

func writeMalformedAssertTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{{Name: "bad_leaf"}})
	content := "---\nlabel: [broken\n---\n\n## Expected\n\n" + goFence() + "func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n" + endFence()
	assertPath := filepath.Join(root, "bad_leaf", "ASSERT.md")
	if err := os.WriteFile(assertPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeUnlabeledTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeaf(t, root, "plain_leaf")
	return root
}

func treeInWorkDir(t *testing.T, name string, includeFast bool, label, explanation string) string {
	t.Helper()
	workDir := t.TempDir()
	treeDir := filepath.Join(workDir, name)
	if err := os.MkdirAll(treeDir, 0755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteFile(t, treeDir, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	if includeFast {
		writeMinimalLeaf(t, treeDir, "fast_leaf")
	}
	writeMinimalLeaf(t, treeDir, "labeled_leaf")
	writeLabeledAssert(t, treeDir, "labeled_leaf", label, explanation)
	return workDir
}

func skipBlock(stdout string) string {
	start := strings.Index(stdout, "SKIPPED ")
	if start < 0 {
		return ""
	}
	rest := stdout[start:]
	end := strings.Index(rest, "\n\n")
	if end < 0 {
		return strings.TrimRight(rest, "\n")
	}
	return rest[:end]
}

func wantSkipBlock(treeRoot, leafRel, label, explanation string) string {
	block := fmt.Sprintf("SKIPPED 1 TESTS\n  %s\n    label: %s", libdocbuild.SkippedDisplayPath(treeRoot, leafRel), label)
	if explanation != "" {
		block += fmt.Sprintf("\n    explanation: %s", explanation)
	}
	return block
}

func assertSkipBlockExact(t *testing.T, stdout, treeRoot, leafRel, label, explanation string) {
	t.Helper()
	want := wantSkipBlock(treeRoot, leafRel, label, explanation)
	got := skipBlock(stdout)
	if got != want {
		t.Fatalf("skip block mismatch\nwant:\n%s\ngot:\n%s\nstdout:\n%s", want, got, stdout)
	}
}

func assertNoSkipBlock(t *testing.T, stdout string) {
	t.Helper()
	if skipBlock(stdout) != "" {
		t.Fatalf("expected no skip block\nstdout:\n%s", stdout)
	}
}

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func findResultSummary(stdout string) string {
	for _, line := range strings.Split(stdout, "\n") {
		plain := strings.TrimSpace(stripANSI(line))
		if strings.HasPrefix(plain, "PASS (") || strings.HasPrefix(plain, "FAIL (") {
			return line
		}
	}
	return ""
}

func assertResultSummary(t *testing.T, stdout string, passed, total int) {
	t.Helper()
	summary := findResultSummary(stdout)
	if summary == "" {
		t.Fatalf("missing PASS/FAIL summary\nstdout:\n%s", stdout)
	}
	plain := stripANSI(strings.TrimSpace(summary))
	wantPrefix := fmt.Sprintf("PASS (%d/%d) in ", passed, total)
	if !strings.HasPrefix(plain, wantPrefix) {
		t.Fatalf("summary prefix mismatch\nwant prefix: %q\ngot: %q", wantPrefix, plain)
	}
	if !finalSummaryPlainRe.MatchString(plain) {
		t.Fatalf("summary must end with duration: %q", plain)
	}
}

func assertNoResultSummary(t *testing.T, stdout string) {
	t.Helper()
	if findResultSummary(stdout) != "" {
		t.Fatalf("expected no PASS/FAIL summary\nstdout:\n%s", stdout)
	}
}

func readAssertFile(t *testing.T, leafDir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(leafDir, "ASSERT.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeGroupingLabeledTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	var groupSetup strings.Builder
	groupSetup.WriteString("# Scenario\n\n**Feature**: e2e grouping fixture\n\n")
	groupSetup.WriteString(bt(3) + "\ne2e grouping\n" + bt(3) + "\n\n## Steps\n1. e2e grouping node\n\n")
	groupSetup.WriteString(goFence())
	groupSetup.WriteString("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
	groupSetup.WriteString(endFence())
	groupSetup.WriteString("\n")
	testtree.WriteFile(t, root, "e2e/SETUP.md", groupSetup.String())
	writeMinimalLeaf(t, root, "e2e/fast_child")
	writeMinimalLeaf(t, root, "e2e/labeled_child")
	writeLabeledAssert(t, root, "e2e/labeled_child", "ui-automation", "grouping skip")
	return root
}

func writeMultiArgModTree(t *testing.T) (workDir string, explicitLeaf string) {
	t.Helper()
	workDir = t.TempDir()
	modDir := filepath.Join(workDir, "mod")
	testtree.WriteFile(t, modDir, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeaf(t, modDir, "fast_leaf")
	writeMinimalLeaf(t, modDir, "skip_labeled")
	writeLabeledAssert(t, modDir, "skip_labeled", "ui-automation", "discovery skip")
	writeMinimalLeaf(t, modDir, "explicit_labeled")
	writeLabeledAssert(t, modDir, "explicit_labeled", "ui-automation", "explicit run")
	explicitLeaf = filepath.Join(modDir, "explicit_labeled")
	return workDir, explicitLeaf
}

func wantSkipBlockMulti(count int, entries []string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("SKIPPED %d TESTS", count))
	for _, entry := range entries {
		b.WriteString("\n")
		b.WriteString(entry)
	}
	return b.String()
}

func wantSkipEntry(displayPath, label, explanation string) string {
	entry := fmt.Sprintf("  %s\n    label: %s", displayPath, label)
	if explanation != "" {
		entry += fmt.Sprintf("\n    explanation: %s", explanation)
	}
	return entry
}

func assertSkipBlockExactMulti(t *testing.T, stdout string, want string) {
	t.Helper()
	got := skipBlock(stdout)
	if got != want {
		t.Fatalf("skip block mismatch\nwant:\n%s\ngot:\n%s\nstdout:\n%s", want, got, stdout)
	}
}

func wantIdempotentLabelWarning(leafDir, label string) string {
	return fmt.Sprintf("warning: label %q already set on %s\n", label, leafDir)
}

func assertStderrExact(t *testing.T, stderr, want string) {
	t.Helper()
	if stderr != want {
		t.Fatalf("stderr mismatch\nwant:\n%q\ngot:\n%q", want, stderr)
	}
}
```