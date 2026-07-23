package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Nested doctest test via RunWithWriter must put progress/summary in the buffer
// (not leak to the process stdout). Regression: L2 dual-mode harnesses saw empty
// Response.Stdout and outer suites dumped nested logs as fail detail.
func TestRunWithWriterCapturesNestedTestSummary(t *testing.T) {
	root := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{
		{Name: "a"},
		{Name: "b"},
	})
	// Isolate gen so we do not touch the warm mapping-gen of the host suite.
	gen := filepath.Join(t.TempDir(), "gen")
	var buf bytes.Buffer
	err := RunWithWriter(&buf, []string{"test", "--no-color", "--gen-dir", gen, "-count=1", root})
	out := buf.String()
	if err != nil {
		t.Fatalf("RunWithWriter: %v\nbuffer:\n%s", err, out)
	}
	if !strings.Contains(out, "2 Run") || !strings.Contains(out, "2 Pass") {
		t.Fatalf("expected captured progress (2 Run, 2 Pass, …), got:\n%s", out)
	}
	if !strings.Contains(out, "PASS (2/2)") {
		t.Fatalf("expected PASS (2/2) in capture buffer, got:\n%s", out)
	}
	// Sanity: buffer is the capture sink (non-empty is enough; we do not
	// reassign os.Stdout so we cannot assert process silence without a pipe).
	if out == "" {
		t.Fatal("empty capture")
	}
	_ = os.DevNull
}

func TestRunWithWriterCapturesNoTestsOnStderrPath(t *testing.T) {
	empty := t.TempDir()
	var buf bytes.Buffer
	err := RunWithWriter(&buf, []string{"test", empty})
	if err != nil {
		t.Fatalf("empty dir should soft-succeed, got %v", err)
	}
	if !strings.Contains(buf.String(), "no tests") {
		t.Fatalf("expected \"no tests\" in capture, got %q", buf.String())
	}
}

func TestRunWithWritersSplitsStdoutStderr(t *testing.T) {
	root := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{{Name: "a"}})
	gen := filepath.Join(t.TempDir(), "gen")
	var stdout, stderr bytes.Buffer
	err := RunWithWriters(&stdout, &stderr, []string{"test", "-v", "--no-color", "--gen-dir", gen, "-count=1", root})
	if err != nil {
		t.Fatalf("RunWithWriters: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	// cd … go test -v is prepare/display on stderr; PASS/progress on stdout.
	if !strings.Contains(stderr.String(), "go test -v") {
		t.Fatalf("expected go test -v on stderr, got:\n%s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "PASS (1/1)") {
		t.Fatalf("expected PASS on stdout, got:\n%s", stdout.String())
	}
}
