package metrics

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParentLeaf(t *testing.T) {
	ClearParentLeaf()
	if ParentLeaf() != "" {
		t.Fatalf("want empty, got %q", ParentLeaf())
	}
	SetParentLeaf("recording/x")
	if ParentLeaf() != "recording/x" {
		t.Fatalf("got %q", ParentLeaf())
	}
	ClearParentLeaf()
	if ParentLeaf() != "" {
		t.Fatalf("want empty after clear, got %q", ParentLeaf())
	}
	// Process env is ignored (parallel-safe): ParentLeaf is Options / SetParentLeaf
	// only — never reads EnvMetricsParentLeaf. (Do not t.Setenv in unit tests.)
	if ParentLeaf() != "" {
		t.Fatalf("env must not affect ParentLeaf, got %q", ParentLeaf())
	}
	_ = EnvMetricsParentLeaf // keep deprecated const referenced for API stability
}

func TestAppendAndReadNestSink(t *testing.T) {
	dir := t.TempDir()
	sink := filepath.Join(dir, "nest.jsonl")
	if err := AppendNestPhase(sink, "go_test", "recording/a", "/tmp/tree", 1_000_000, map[string]any{"cases": 1}); err != nil {
		t.Fatal(err)
	}
	if err := AppendNestPhase(sink, "generate", "recording/a", "/tmp/tree", 200_000, nil); err != nil {
		t.Fatal(err)
	}
	events, err := ReadNestSinkEvents(sink)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("events=%d", len(events))
	}
	if events[0]["phase"] != "go_test" {
		t.Fatalf("first phase=%v", events[0]["phase"])
	}
	// missing file
	miss, err := ReadNestSinkEvents(filepath.Join(dir, "nope.jsonl"))
	if err != nil || miss != nil {
		t.Fatalf("missing: %v %v", miss, err)
	}
	_ = os.Remove(sink)
}
