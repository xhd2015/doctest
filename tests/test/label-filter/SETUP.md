# Scenario

**Feature**: `doctest test --label` selects labeled leaves by boolean expression

```
# build doctest binary once per run
Setup -> build ./cmd/doctest -> req.Bin

# subprocess exercises CLI contract
req.Bin <args> -> stdout/stderr/exit captured in Response
```

## Preconditions

- Fresh doctest binary built from module source in this root Setup.

## Steps

1. Build a temp fixture mod tree when the leaf needs integration coverage.
2. Configure `req.Args` and optional `req.WorkDir`.
3. Assert observable stdout/stderr blocks.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
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
	tmp := t.TempDir()
	doctestBin := filepath.Join(tmp, "doctest")
	buildDir := filepath.Join(DOCTEST_ROOT, "..", "..", "..")
	buildArgs := []string{"build", "-o", doctestBin}
	if libdocbuild.NeedsBuildVCSFlag(buildDir) {
		buildArgs = append(buildArgs, "-buildvcs=false")
	}
	buildArgs = append(buildArgs, "./cmd/doctest")
	build := exec.Command("go", buildArgs...)
	build.Dir = buildDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build doctest: %v\n%s", err, string(out))
	}
	req.Bin = doctestBin
	return nil
}

func writeMinimalLeafSetupAssert(t *testing.T, root, leafRel string) {
	t.Helper()
	var setup strings.Builder
	setup.WriteString("# Scenario\n\n**Feature**: label-filter fixture leaf\n\n")
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

// writeLabelFilterMod builds the standard five-leaf mod from the requirement.
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

func wantLabelFilterSkipEntry(treeRoot, leafRel, label, explanation string, unlabeled bool) string {
	display := libdocbuild.SkippedDisplayPath(treeRoot, leafRel)
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(display)
	if !unlabeled {
		b.WriteString("\n    label: ")
		b.WriteString(label)
	}
	if explanation != "" {
		b.WriteString("\n    explanation: ")
		b.WriteString(explanation)
	}
	b.WriteString("\n    reason: label filter")
	return b.String()
}

func wantLabelFilterSkipBlockMulti(count int, entries ...string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("SKIPPED %d TESTS", count))
	for _, e := range entries {
		b.WriteString("\n")
		b.WriteString(e)
	}
	return b.String()
}

func assertSkipBlockExact(t *testing.T, stdout, want string) {
	t.Helper()
	got := skipBlock(stdout)
	if got != want {
		t.Fatalf("skip block mismatch\nwant:\n%s\ngot:\n%s\nstdout:\n%s", want, got, stdout)
	}
}

func assertSkipBlockContainsReason(t *testing.T, stdout string) {
	t.Helper()
	block := skipBlock(stdout)
	if block == "" {
		t.Fatalf("expected skip block with reason\nstdout:\n%s", stdout)
	}
	if !strings.Contains(block, "reason: label filter") {
		t.Fatalf("skip block missing reason line\nblock:\n%s", block)
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

func assertStderrContains(t *testing.T, stderr, substr string) {
	t.Helper()
	if !strings.Contains(stderr, substr) {
		t.Fatalf("stderr missing %q:\n%s", substr, stderr)
	}
}

func sortedLeafRels(names ...string) []string {
	out := append([]string(nil), names...)
	sort.Strings(out)
	return out
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, string(out))
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "label-filter@test.com")
	runGit(t, dir, "config", "user.name", "Label Filter Test")
}

func gitAddCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-m", msg)
}

func writeLabelFilterModInDir(t *testing.T, modDir string) {
	t.Helper()
	if err := os.MkdirAll(modDir, 0755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteFile(t, modDir, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	writeMinimalLeafSetupAssert(t, modDir, "fast")
	writeMinimalLeafSetupAssert(t, modDir, "slow")
	writeLabeledAssert(t, modDir, "slow", "slow", "slow profile")
	writeMinimalLeafSetupAssert(t, modDir, "heavy")
	writeLabeledAssert(t, modDir, "heavy", "heavy", "heavy profile")
	writeMinimalLeafSetupAssert(t, modDir, "both")
	writeLabeledAssert(t, modDir, "both", "slow, ui-automation", "slow ui combo")
}

func createLabelFilterGitMod(t *testing.T) (string, string) {
	t.Helper()
	repoDir := t.TempDir()
	modDir := filepath.Join(repoDir, "mod")
	writeLabelFilterModInDir(t, modDir)
	initGitRepo(t, repoDir)
	gitAddCommitAll(t, repoDir, "baseline label-filter mod")
	return repoDir, modDir
}
```