package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTreeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func setupDoc(code string) string {
	code = trimDocCode(code)
	if !containsDocFunc(code, "func Setup") {
		setupLine := "func Setup(t *testing.T, req *Request) error { _ = req; return nil }"
		code = injectSetupFunc(code, setupLine)
	}
	return "# Setup\n\nAny section names are allowed.\n\n```go\n" + code + "\n```\n"
}

func assertDoc(code string) string {
	return "# Assert\n\nAny section names are allowed.\n\n```go\n" + trimDocCode(code) + "\n```\n"
}

func trimDocCode(code string) string {
	for len(code) > 0 && (code[0] == '\n' || code[0] == '\r') {
		code = code[1:]
	}
	for len(code) > 0 && (code[len(code)-1] == '\n' || code[len(code)-1] == '\r') {
		code = code[:len(code)-1]
	}
	return code
}

func containsDocFunc(code, fn string) bool {
	return bytes.Contains([]byte(code), []byte(fn))
}

func injectSetupFunc(code, setupLine string) string {
	idx := importEnd(code)
	if idx >= 0 {
		return code[:idx] + "\n" + setupLine + "\n" + code[idx:]
	}
	return setupLine + "\n" + code
}

func importEnd(code string) int {
	for _, delim := range []string{"\")\n", "\"\n"} {
		if idx := indexAfter(code, delim); idx >= 0 {
			return idx
		}
	}
	return -1
}

func indexAfter(s, substr string) int {
	idx := bytes.Index([]byte(s), []byte(substr))
	if idx < 0 {
		return -1
	}
	return idx + len(substr)
}

func createValidTestTree(t *testing.T, root string) {
	t.Helper()
	writeTreeFile(t, root, "DOCTEST.md", "# tree")
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))
}

func TestBuildArgsWithValidTree(t *testing.T) {
	root := t.TempDir()
	createValidTestTree(t, root)
	genDir := filepath.Join(t.TempDir(), "gen")

	args := []string{"--rm", "--gen-dir", genDir, root}
	if err := BuildArgs(args); err != nil {
		t.Fatalf("BuildArgs: %v", err)
	}

	entries, err := os.ReadDir(genDir)
	if err != nil {
		t.Fatalf("read gen dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected generated files")
	}
}

func TestBuildArgsMissingDir(t *testing.T) {
	err := BuildArgs([]string{})
	if err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestWithArgsMissingDir(t *testing.T) {
	err := Test([]string{})
	if err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestBuildArgsVerbose(t *testing.T) {
	root := t.TempDir()
	createValidTestTree(t, root)
	genDir := filepath.Join(t.TempDir(), "gen")

	args := []string{"--verbose", "--rm", "--gen-dir", genDir, root}
	if err := BuildArgs(args); err != nil {
		t.Fatalf("BuildArgs: %v", err)
	}
}

func TestWithArgsWithValidTree(t *testing.T) {
	root := t.TempDir()
	createValidTestTree(t, root)

	args := []string{"--rm", root}
	if err := Test(args); err != nil {
		t.Fatalf("Test: %v", err)
	}
}

func TestParseBuildOptionsRemoveTempDefault(t *testing.T) {
	opts, _, err := parseBuildOptions([]string{"somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.RemoveTemp {
		t.Fatal("expected RemoveTemp=false by default, got true")
	}
}

func TestParseBuildOptionsRemoveTempFlag(t *testing.T) {
	opts, _, err := parseBuildOptions([]string{"--rm", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.RemoveTemp {
		t.Fatal("expected RemoveTemp=true with --rm flag, got false")
	}
}

func TestParseTestOptionsRemoveTempDefault(t *testing.T) {
	opts, _, err := parseTestOptions([]string{"somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.RemoveTemp {
		t.Fatal("expected RemoveTemp=false by default, got true")
	}
}

func TestParseTestOptionsRemoveTempFlag(t *testing.T) {
	opts, _, err := parseTestOptions([]string{"--rm", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.RemoveTemp {
		t.Fatal("expected RemoveTemp=true with --rm flag, got false")
	}
}
