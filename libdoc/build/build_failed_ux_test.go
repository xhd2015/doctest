package build

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// TestBuildFailedUX_noPhantomRunOrCached asserts that when the generated suite
// fails to compile, quiet progress does not invent "1 Run" from package fail
// or "N Cached" from a pre-planned leaf-cache skip list, and the footer uses
// planned denom: FAIL (build failed; 0/N executed).
func TestBuildFailedUX_noPhantomRunOrCached(t *testing.T) {
	root := t.TempDir()
	// Run body references undefined symbol → suite package [build failed].
	brokenRun := `import "testing"

type Request struct{}
type Response struct{}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = definitelyNotDefined
	return &Response{}, nil
}`
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(brokenRun))
	// Two leaves so Planned=2 (not 1 — guards against package-fail denom).
	for _, name := range []string{"leaf_a", "leaf_b"} {
		testtree.WriteFile(t, root, name+"/SETUP.md", "## Steps\n1. x\n\n```go\nimport \"testing\"\n\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }\n```\n")
		testtree.WriteFile(t, root, name+"/ASSERT.md", "## Expected\n- ok\n\n```go\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n```\n")
	}

	genDir := filepath.Join(t.TempDir(), "gen")
	var stdout, stderr bytes.Buffer
	opts := core.Options{
		GenDir:                genDir,
		Count:                 1, // disable leaf-cache skip so Cached=0 is unambiguous
		NoLeafCache:           true,
		SuppressResultSummary: false,
		Stdout:                &stdout,
		Stderr:                &stderr,
		Color:                 core.ColorNever,
	}
	stats, err := TestWithStats(root, opts)
	if err == nil {
		t.Fatalf("expected go test error on build failed, got nil\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
	out := stdout.String()
	// Package line may include [build failed].
	if !strings.Contains(out, "[build failed]") && !strings.Contains(stderr.String(), "[build failed]") && !stats.BuildFailed {
		t.Fatalf("expected build failed signal\nstdout:\n%s\nstderr:\n%s\nstats=%+v", out, stderr.String(), stats)
	}
	if !stats.BuildFailed {
		t.Fatalf("stats.BuildFailed=false; planned=%d total=%d passed=%d\n%s", stats.Planned, stats.Total, stats.Passed, out)
	}
	if stats.Planned != 2 {
		t.Fatalf("planned=%d want 2", stats.Planned)
	}
	if stats.Total != 0 || stats.Passed != 0 {
		t.Fatalf("executed must be 0/0, got pass=%d total=%d", stats.Passed, stats.Total)
	}
	// Progress summary: no phantom 1 Run / N Cached.
	if !strings.Contains(out, "0 Run") || !strings.Contains(out, "0 Cached") {
		t.Fatalf("want (0 Run, …, 0 Cached), got:\n%s", out)
	}
	if strings.Contains(out, "1 Run") || strings.Contains(out, "1 Fail") {
		t.Fatalf("must not inflate package fail as leaf Run/Fail:\n%s", out)
	}
	// Footer uses planned denom, not FAIL (0/1).
	if !strings.Contains(out, "FAIL (build failed; 0/2 executed)") {
		t.Fatalf("want build-failed footer with planned=2:\n%s", out)
	}
	if strings.Contains(out, "FAIL (0/1)") {
		t.Fatalf("must not use package-fail denom FAIL (0/1):\n%s", out)
	}
}

// TestBuildFailedUX_workspaceWithWarmCachePlan ensures Cached stays 0 even when
// a leaf-cache skip plan would have been non-empty (skip list ignored).
func TestBuildFailedUX_workspaceWithWarmCachePlan(t *testing.T) {
	// Unit-level: applyGoTestLeafStats already covers skipPaths; this checks
	// formatSummary composition with planned multi-leaf.
	stats := TestRunStats{Total: 5, Planned: 5}
	result := goTestJSONResult{buildFailed: true}
	skipPaths := []string{"l1", "l2", "l3", "l4", "l5"}
	applyGoTestLeafStats(&stats, &result, 5, skipPaths, core.Options{})
	if result.cachedCount != 0 || !stats.BuildFailed {
		t.Fatalf("cached=%d buildFailed=%v", result.cachedCount, stats.BuildFailed)
	}
	sum := formatSummary(colorStyle{}, result.passCount+result.failCount, result.passCount, result.failCount, result.cachedCount, 0)
	foot := formatBuildFailedResultSummary(colorStyle{}, stats.Planned, 0)
	if !strings.Contains(sum, "0 Run") || !strings.Contains(sum, "0 Cached") {
		t.Fatalf("sum=%q", sum)
	}
	if !strings.HasPrefix(foot, "FAIL (build failed; 0/5 executed)") {
		t.Fatalf("foot=%q", foot)
	}
}
