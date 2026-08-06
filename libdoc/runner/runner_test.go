package runner

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
		setupLine := "func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }"
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

func doctestDoc(code string) string {
	return "# Tests\n\n## Version\n0.0.2\n\n## DSN (Domain Specific Notion)\n\n### Participants\n- **system** — under test.\n\n### Behaviors\n- **run** — executes.\n\n```go\n" + trimDocCode(code) + "\n```\n"
}

func createValidTestTree(t *testing.T, root string) {
	t.Helper()
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
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

func TestParseTestOptionsTimeoutFlag(t *testing.T) {
	opts, remain, err := parseTestOptions([]string{"-timeout", "45s", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Timeout == nil || *opts.Timeout != 45*time.Second {
		t.Fatalf("expected Timeout=45s, got %v", opts.Timeout)
	}
	if len(remain) != 1 || remain[0] != "somedir" {
		t.Fatalf("expected remainArgs [somedir], got %v", remain)
	}
}

func TestParseTestOptionsTimeoutLongAlias(t *testing.T) {
	opts, _, err := parseTestOptions([]string{"--timeout=45s", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Timeout == nil || *opts.Timeout != 45*time.Second {
		t.Fatalf("expected --timeout alias Timeout=45s, got %v", opts.Timeout)
	}
}

func TestParseTestOptionsTimeoutZeroDisables(t *testing.T) {
	opts, _, err := parseTestOptions([]string{"-timeout", "0", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Timeout == nil {
		t.Fatal("expected Timeout non-nil for -timeout=0 (disable), got nil")
	}
	if *opts.Timeout != 0 {
		t.Fatalf("expected Timeout=0, got %v", *opts.Timeout)
	}
}

func TestParseTestOptionsTimeoutOmitted(t *testing.T) {
	opts, _, err := parseTestOptions([]string{"somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Timeout != nil {
		t.Fatalf("expected Timeout nil when -timeout omitted, got %v", opts.Timeout)
	}
}

func TestParseTestOptionsTimeoutInvalid(t *testing.T) {
	_, _, err := parseTestOptions([]string{"-timeout", "bogus", "somedir"})
	if err == nil {
		t.Fatal("expected error for invalid -timeout value")
	}
}

func TestParseTestOptionsRace(t *testing.T) {
	opts, remain, err := parseTestOptions([]string{"-race", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Race {
		t.Fatal("expected Race true")
	}
	if len(remain) != 1 || remain[0] != "somedir" {
		t.Fatalf("remain=%v", remain)
	}
	optsOff, _, err := parseTestOptions([]string{"somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if optsOff.Race {
		t.Fatal("expected Race false by default")
	}
}

func TestParseTestOptionsGoTestAllowlist(t *testing.T) {
	opts, remain, err := parseTestOptions([]string{
		"-covermode", "atomic",
		"-coverpkg", "example.com/mod/...",
		"-short",
		"-failfast",
		"-parallel", "4",
		"-shuffle", "on",
		"-tags", "integration",
		"-gcflags", "all=-N",
		"-ldflags", "-X=main.v=1",
		"somedir",
	})
	if err != nil {
		t.Fatal(err)
	}
	if opts.CoverMode != "atomic" || opts.CoverPkg != "example.com/mod/..." {
		t.Fatalf("cover fields: mode=%q pkg=%q", opts.CoverMode, opts.CoverPkg)
	}
	if !opts.Short || !opts.FailFast {
		t.Fatalf("short=%v failfast=%v", opts.Short, opts.FailFast)
	}
	if opts.Parallel == nil || *opts.Parallel != 4 {
		t.Fatalf("parallel=%v", opts.Parallel)
	}
	if opts.Shuffle != "on" || opts.Tags != "integration" {
		t.Fatalf("shuffle=%q tags=%q", opts.Shuffle, opts.Tags)
	}
	if opts.Gcflags != "all=-N" || opts.Ldflags != "-X=main.v=1" {
		t.Fatalf("gcflags=%q ldflags=%q", opts.Gcflags, opts.Ldflags)
	}
	if len(remain) != 1 || remain[0] != "somedir" {
		t.Fatalf("remain=%v", remain)
	}
}

func TestParseTestOptionsRejectNameFilters(t *testing.T) {
	for _, flag := range []string{"-run", "-skip", "-bench", "-benchtime", "-benchmem", "-fuzz", "-run=TestFoo", "--run"} {
		_, _, err := parseTestOptions([]string{flag, "somedir"})
		if err == nil {
			t.Fatalf("%s: expected reject", flag)
		}
		msg := err.Error()
		if !strings.Contains(msg, "not supported") || !strings.Contains(msg, "--label") {
			t.Fatalf("%s: unexpected err: %v", flag, err)
		}
	}
}

func TestParseTestOptionsLabelAll(t *testing.T) {
	opts, remain, err := parseTestOptions([]string{"--label-all", "somedir"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.LabelAll {
		t.Fatal("expected LabelAll=true")
	}
	if len(remain) != 1 || remain[0] != "somedir" {
		t.Fatalf("remain=%v", remain)
	}
}

func TestParseTestOptionsLabelAllConflictsWithLabel(t *testing.T) {
	_, _, err := parseTestOptions([]string{"--label-all", "--label", "e2e", "somedir"})
	if err == nil {
		t.Fatal("expected mutual exclusion error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("err=%v", err)
	}
}

func TestFormatClassifiedErrors(t *testing.T) {
	if err := formatClassifiedErrors(nil, nil); err != nil {
		t.Fatalf("empty: %v", err)
	}
	prepOnly := formatClassifiedErrors([]string{"/t/bad: go mod tidy: exit 1"}, nil)
	if prepOnly == nil || !strings.HasPrefix(prepOnly.Error(), "prepare failed:\n") {
		t.Fatalf("prepare-only: %v", prepOnly)
	}
	if !strings.Contains(prepOnly.Error(), "/t/bad:") {
		t.Fatalf("prepare-only body: %v", prepOnly)
	}
	runOnly := formatClassifiedErrors(nil, []string{"workspace: build failed"})
	if runOnly == nil || !strings.HasPrefix(runOnly.Error(), "test failures:\n") {
		t.Fatalf("run-only: %v", runOnly)
	}
	mixed := formatClassifiedErrors([]string{"a: prep"}, []string{"b: run"})
	if mixed == nil || !strings.HasPrefix(mixed.Error(), "errors:\n") {
		t.Fatalf("mixed: %v", mixed)
	}
	if !strings.Contains(mixed.Error(), "a: prep") || !strings.Contains(mixed.Error(), "b: run") {
		t.Fatalf("mixed body: %v", mixed)
	}
}

// TestDotDotDotPrepareFailNoPASS: one tree fails prepare, sibling runs and
// passes — overall must not print PASS (honest FAIL) and error is prepare failed.
func TestDotDotDotPrepareFailNoPASS(t *testing.T) {
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/prepfail\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	good := filepath.Join(mod, "good")
	createValidTestTree(t, good)

	bad := filepath.Join(mod, "bad")
	// Invalid Go in DOCTEST.md so generate/prepare fails.
	writeTreeFile(t, bad, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil
// missing closing brace — syntax error
`))
	writeTreeFile(t, bad, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, bad, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "gen")
	// Absolute mod/... — no process os.Chdir (Parallel-safe).
	pattern := filepath.Join(mod, "...")

	// Inject opts.Stdout — never reassign os.Stdout.
	var outBuf bytes.Buffer
	opts, remain, err := parseTestOptions([]string{"--gen-dir", genDir, "--no-color", pattern})
	if err != nil {
		t.Fatal(err)
	}
	opts.Stdout = &outBuf
	opts.Stderr = io.Discard
	runErr := runTest(opts, remain)
	out := outBuf.String()

	if runErr == nil {
		t.Fatal("expected non-nil error when one tree fails prepare")
	}
	if !strings.Contains(runErr.Error(), "prepare failed:") {
		t.Fatalf("expected prepare failed label, got: %v", runErr)
	}
	if strings.Contains(out, "PASS (") {
		t.Fatalf("must not print PASS when prepare failed:\n%s", out)
	}
	// Good tree ran one case — summary must be honest FAIL, not silent.
	if !strings.Contains(out, "FAIL (") {
		t.Fatalf("expected FAIL summary when survivors ran:\nstdout=%q\nerr=%v", out, runErr)
	}
}

// TestExpandTestArgsMidPathSubDir checks that path/... under a mid branch sets
// SubDir to the mid path (not only the full DOCTEST tree root).
func TestExpandTestArgsMidPathSubDir(t *testing.T) {
	// tree/DOCTEST.md is the suite root; mid/ is a sub-branch with leaves.
	// expandTestArgs("tree/mid/...") must filter under mid, not the whole tree.
	base := t.TempDir()
	tree := filepath.Join(base, "tree")
	mid := filepath.Join(tree, "mid")
	writeTreeFile(t, tree, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, mid, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, mid, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))
	writeTreeFile(t, tree, "sibling/SETUP.md", setupDoc(`
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, tree, "sibling/ASSERT.md", assertDoc(`
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}
`))

	targets, err := expandTestArgs([]string{mid + "/..."})
	if err != nil {
		t.Fatalf("expandTestArgs: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("want 1 target, got %d: %+v", len(targets), targets)
	}
	midClean := filepath.Clean(mid)
	treeClean := filepath.Clean(tree)
	gotRoot := filepath.Clean(targets[0].Root)
	gotSub := filepath.Clean(targets[0].SubDir)
	if gotRoot != treeClean {
		t.Errorf("Root=%q want %q", gotRoot, treeClean)
	}
	if gotSub != midClean {
		t.Errorf("SubDir=%q want mid %q (not full tree only)", gotSub, midClean)
	}
}
