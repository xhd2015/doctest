package core

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverLightSkipsDeepParseOfHeavy(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	// Unlabeled good leaf
	writeTreeFile(t, root, "fast/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "fast/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))
	// Heavy leaf with BROKEN setup body — would fail full discover hydrate of that leaf
	writeTreeFile(t, root, "slow/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "slow/ASSERT.md", "---\nlabel: heavy\n---\n\n## Expected\n\n```go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tnot valid go (((\n}\n```\n")

	light, err := DiscoverTreeCasesLight(root)
	if err != nil {
		t.Fatalf("light: %v", err)
	}
	if len(light) != 2 {
		t.Fatalf("light cases=%d", len(light))
	}
	run, skipped := FilterCasesByLabel(light, Options{})
	if len(run) != 1 || len(skipped) != 1 {
		t.Fatalf("run=%d skipped=%d", len(run), len(skipped))
	}
	if skipped[0].Labels[0] != "heavy" {
		t.Fatalf("skipped labels=%v", skipped[0].Labels)
	}
	// Hydrate only run set — must succeed despite broken heavy ASSERT
	hydrated, err := HydrateTreeCases(root, run)
	if err != nil {
		t.Fatalf("hydrate run set: %v", err)
	}
	if len(hydrated) != 1 || hydrated[0].Path != "fast" {
		t.Fatalf("hydrated=%v", hydrated)
	}
	if hydrated[0].AssertFile.GoBlock.Assert == nil {
		t.Fatal("expected deep assert")
	}
	// Full discover still fails on heavy
	if _, err := DiscoverTreeCases(root); err == nil {
		t.Fatal("full discover should fail on broken heavy assert")
	}
}

func TestHydrateBrokenSelectedFails(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "bad/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "bad/ASSERT.md", "## Expected\n\n```go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tnot valid (((\n}\n```\n")

	light, err := DiscoverTreeCasesLight(root)
	if err != nil {
		t.Fatal(err)
	}
	run, _ := FilterCasesByLabel(light, Options{})
	_, err = HydrateTreeCases(root, run)
	if err == nil {
		t.Fatal("expected hydrate error")
	}
	if !strings.Contains(err.Error(), "ASSERT") && !strings.Contains(err.Error(), "bad") {
		// parse errors mention path
		t.Logf("err=%v", err)
	}
	_ = filepath.Separator
}
