package core

import (
	"strings"
	"testing"
)

func TestParseAssertFrontmatterMissing(t *testing.T) {
	fm, body, err := ParseAssertFrontmatter("ASSERT.md", "## Expected\n\n```go\nfunc Assert() {}\n```\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(fm.Labels) != 0 || fm.Explanation != "" {
		t.Fatalf("fm = %#v, want empty", fm)
	}
	if !strings.HasPrefix(body, "## Expected") {
		t.Fatalf("body = %q", body)
	}
}

func TestParseAssertFrontmatterLabelsAndExplanation(t *testing.T) {
	content := "---\nlabel: human-guided-ui-test, debug\nexplanation: used for debugging\n---\n\n## Expected\n"
	fm, body, err := ParseAssertFrontmatter("ASSERT.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(fm.Labels, ",") != "human-guided-ui-test,debug" {
		t.Fatalf("labels = %#v", fm.Labels)
	}
	if fm.Explanation != "used for debugging" {
		t.Fatalf("explanation = %q", fm.Explanation)
	}
	if !strings.HasPrefix(body, "## Expected") {
		t.Fatalf("body = %q", body)
	}
}

func TestParseAssertFrontmatterExplanationOnly(t *testing.T) {
	content := "---\nexplanation: note only\n---\n\n## Expected\n"
	fm, _, err := ParseAssertFrontmatter("ASSERT.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(fm.Labels) != 0 || fm.Explanation != "note only" {
		t.Fatalf("fm = %#v", fm)
	}
}

func TestParseAssertFrontmatterMalformedYAML(t *testing.T) {
	content := "---\nlabel: [unterminated\n---\n\n## Expected\n"
	_, _, err := ParseAssertFrontmatter("ASSERT.md", content)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAssertFrontmatterMissingClose(t *testing.T) {
	content := "---\nlabel: x\n\n## Expected\n"
	_, _, err := ParseAssertFrontmatter("ASSERT.md", content)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestPartitionLabeledCasesDiscoverySkips(t *testing.T) {
	cases := []TreeCase{
		{Name: "fast", Path: "fast"},
		{Name: "slow", Path: "slow", Labels: []string{"ui"}},
	}
	run, skipped := PartitionLabeledCases(cases, true)
	if len(run) != 1 || run[0].Path != "fast" {
		t.Fatalf("run = %#v", run)
	}
	if len(skipped) != 1 || skipped[0].Path != "slow" {
		t.Fatalf("skipped = %#v", skipped)
	}
}

func TestPartitionLabeledCasesExplicitLeafRunsAll(t *testing.T) {
	cases := []TreeCase{{Name: "slow", Path: "slow", Labels: []string{"ui"}}}
	run, skipped := PartitionLabeledCases(cases, false)
	if len(run) != 1 || len(skipped) != 0 {
		t.Fatalf("run=%#v skipped=%#v", run, skipped)
	}
}