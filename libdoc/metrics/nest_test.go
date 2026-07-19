package metrics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParentLeaf(t *testing.T) {
	ClearParentLeaf()
	_ = os.Unsetenv(EnvMetricsParentLeaf)
	if ParentLeaf() != "" {
		t.Fatal("want empty")
	}
	SetParentLeaf("recording/x")
	if ParentLeaf() != "recording/x" {
		t.Fatalf("got %q", ParentLeaf())
	}
	ClearParentLeaf()
	if ParentLeaf() != "" {
		t.Fatal("want empty after clear")
	}
	t.Setenv(EnvMetricsParentLeaf, "via-env")
	if ParentLeaf() != "via-env" {
		t.Fatalf("env got %q", ParentLeaf())
	}
}

func TestAppendAndReadNestSink(t *testing.T) {
	dir := t.TempDir()
	sink := filepath.Join(dir, "run.nest")
	if err := AppendNestPhase(sink, "go_test", "recording/a", "/tmp/fix", 1_000_000, map[string]any{"cases": 1}); err != nil {
		t.Fatal(err)
	}
	if err := AppendNestPhase(sink, "generate", "recording/a", "/tmp/fix", 200_000, nil); err != nil {
		t.Fatal(err)
	}
	events, err := ReadNestSinkEvents(sink)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("events=%d", len(events))
	}
	if events[0]["scope"] != "nested" || events[0]["phase"] != "go_test" {
		t.Fatalf("first=%v", events[0])
	}
	if events[0]["parent_leaf"] != "recording/a" {
		t.Fatalf("parent_leaf=%v", events[0]["parent_leaf"])
	}
	// missing file
	miss, err := ReadNestSinkEvents(filepath.Join(dir, "nope"))
	if err != nil || miss != nil {
		t.Fatalf("missing: %v %v", miss, err)
	}
	_ = os.Remove(sink)
}

func TestRankNested(t *testing.T) {
	rows := []PhaseRow{
		{Scope: "tree", Phase: "go_test", ElapsedNs: 9e9},
		{Scope: "nested", Phase: "go_test", ParentLeaf: "recording/b", ElapsedNs: 2e9},
		{Scope: "nested", Phase: "generate", ParentLeaf: "recording/b", ElapsedNs: 1e8},
		{Scope: "nested", Phase: "go_test", ParentLeaf: "recording/a", ElapsedNs: 5e8},
	}
	byLeaf := RankNestedByParentLeaf(rows, 0)
	if len(byLeaf) != 2 || byLeaf[0].ParentLeaf != "recording/b" {
		t.Fatalf("byLeaf=%+v", byLeaf)
	}
	gt := RankNestedGoTest(rows, 0)
	if len(gt) != 2 || gt[0].ParentLeaf != "recording/b" {
		t.Fatalf("gt=%+v", gt)
	}
}
