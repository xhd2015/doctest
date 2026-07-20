package build

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Two sibling Go modules under one ./... pattern run via toplevel/__hub
// (single go test binary), not go.work.
func TestRunWorkspaceMultiModHubSiblingModules(t *testing.T) {
	base := t.TempDir()
	modA := filepath.Join(base, "mod_a")
	modB := filepath.Join(base, "mod_b")
	for _, m := range []string{modA, modB} {
		if err := os.MkdirAll(m, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(m, "go.mod"), []byte("module example.com/"+filepath.Base(m)+"\n\ngo 1.22\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		testtree.WriteMinimalRunnableTree(t, m, []testtree.LeafSpec{
			{Name: "simple", Steps: "x", Expected: "ok"},
		})
	}

	genA := filepath.Join(t.TempDir(), "gen_a")
	genB := filepath.Join(t.TempDir(), "gen_b")
	optsA := core.Options{GenDir: genA, Count: 1, SuppressResultSummary: true, Stderr: ioDiscard{}}
	optsB := core.Options{GenDir: genB, Count: 1, SuppressResultSummary: true, Stderr: ioDiscard{}}

	prepA, err := PrepareTree(modA, optsA)
	if err != nil {
		t.Fatalf("PrepareTree A: %v", err)
	}
	prepB, err := PrepareTree(modB, optsB)
	if err != nil {
		t.Fatalf("PrepareTree B: %v", err)
	}
	if filepath.Clean(prepA.GenRoot) == filepath.Clean(prepB.GenRoot) {
		t.Fatal("expected distinct gen roots for sibling modules")
	}

	var stdout bytes.Buffer
	runOpts := core.Options{Count: 1, SuppressResultSummary: true, Stdout: &stdout, Stderr: ioDiscard{}}
	stats, err := RunWorkspace([]TreePrep{prepA, prepB}, runOpts)
	if err != nil {
		t.Fatalf("RunWorkspace multi-mod hub: %v\nstdout:\n%s", err, stdout.String())
	}
	if stats.Passed != 2 || stats.Total != 2 {
		t.Fatalf("want 2/2, got pass=%d total=%d\n%s", stats.Passed, stats.Total, stdout.String())
	}

	toplevel := pickToplevelGenRoot([]string{prepA.GenRoot, prepB.GenRoot})
	hubGoMod := filepath.Join(toplevel, HubDirName, "go.mod")
	body, err := os.ReadFile(hubGoMod)
	if err != nil {
		t.Fatalf("hub go.mod: %v", err)
	}
	if !strings.Contains(string(body), "module testcase/hub") {
		t.Fatalf("hub go.mod:\n%s", body)
	}
	if _, err := os.Stat(filepath.Join(toplevel, HubDirName, "go.work")); !os.IsNotExist(err) {
		t.Fatalf("go.work must not exist under hub, err=%v", err)
	}
	// Each member has RunAll
	for _, g := range []string{prepA.GenRoot, prepB.GenRoot} {
		// suite/runall.go under tree
		found := false
		_ = filepath.Walk(g, func(p string, info os.FileInfo, err error) error {
			if info != nil && info.Name() == "runall.go" {
				found = true
			}
			return nil
		})
		if !found {
			t.Fatalf("missing runall.go under %s", g)
		}
	}
}

func TestPickToplevelGenRootNested(t *testing.T) {
	parent := "/tmp/gen/parent"
	child := "/tmp/gen/parent/sub"
	got := pickToplevelGenRoot([]string{child, parent})
	if got != parent {
		t.Fatalf("got %q want %q", got, parent)
	}
}
