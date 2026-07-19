# Scenario

**Feature**: `doctest test --cold-cache` runs a reproducible no-warm-cache baseline

```
# cold-cache mode (test command only)
user -> doctest test --cold-cache [ --gen-dir X ] <tiny-tree>
  -> resolve gen root:
       omit --gen-dir  => $CacheHome/doctest/mapping-gen-cold
       --gen-dir X     => allow iff abs(X) not equal/under warm mapping-gen
  -> startup: RemoveAll+MkdirAll chosen gen root (leftover kept after finish)
  -> force -count=1 when count unset; isolate empty GOCACHE for the run
  -> announce cold-cache mode on stderr
  -> generate + go test under cold root

# cache home
DOCTEST_CACHE_HOME sandbox -> CacheHome() so tests never touch developer cache
```

## Preconditions

- Root harness (`tests/SETUP.md`) builds/resolves `req.Bin` via `testbin.Ensure`.
- All leaves set an isolated `DOCTEST_CACHE_HOME` so warm/cold mapping-gen paths
  stay under a per-leaf temp sandbox.
- Fixture trees are minimal (1 leaf) so leaves stay fast.
- `--cold-cache` is **not** implemented yet (classic TDD RED until implementer).

## Steps

1. Raise timeout for generate + `go test` of a tiny fixture.
2. Provide helpers: minimal project tree, cache sandbox, warm/cold home paths,
   marker seed, go-test command line extraction.

## Context

- Warm home: `$CacheHome/doctest/mapping-gen`
- Cold home (auto): `$CacheHome/doctest/mapping-gen-cold`
- Help coverage for the flag lives in `tests/help/test-options` (updated token list).

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"

// st carries per-leaf paths for Assert (reset in each leaf Setup).
var st struct {
	CacheHome string
	WarmHome  string
	ColdHome  string
	GenDir    string // explicit --gen-dir when set
	TestDir   string
	Marker    string // path of pre-seeded marker file
}

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second
	st = struct {
		CacheHome string
		WarmHome  string
		ColdHome  string
		GenDir    string
		TestDir   string
		Marker    string
	}{}
	return nil
}

func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func doctestBody(extraRunCode string) string {
	return "import \"testing\"\n\ntype Request struct{ Args []string; WorkDir string }\ntype Response struct{ ExitCode int; Stdout string; Stderr string }\n\n" + extraRunCode
}

func rootSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafSetupContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafAssertContent() string {
	return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

func createTinyTestTree(dir string) error {
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody("func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"))), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
		return err
	}
	leafDir := filepath.Join(dir, "simple")
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
		return err
	}
	return nil
}

// createTempTestProject returns the doctest tree dir inside a temp module.
func createTempTestProject(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	testDir := filepath.Join(tmp, "mytest")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}
	if err := createTinyTestTree(testDir); err != nil {
		t.Fatalf("create tiny test tree: %v", err)
	}
	return testDir
}

// withCacheSandbox sets DOCTEST_CACHE_HOME on the subprocess env and fills st paths.
func withCacheSandbox(t *testing.T, req *Request) {
	t.Helper()
	cache := t.TempDir()
	st.CacheHome = cache
	st.WarmHome = filepath.Join(cache, "doctest", "mapping-gen")
	st.ColdHome = filepath.Join(cache, "doctest", "mapping-gen-cold")
	req.Env = append(req.Env, "DOCTEST_CACHE_HOME="+cache)
}

func seedMarker(t *testing.T, dir, name string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir for marker: %v", err)
	}
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("marker-before\n"), 0644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	st.Marker = p
	return p
}

// goTestCmdLine returns the first "cd … && go test …" preview line from stderr.
func goTestCmdLine(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "cd ") && strings.Contains(trimmed, " && go test") {
			return trimmed
		}
	}
	return ""
}

func dirHasGoFiles(root string) bool {
	found := false
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".go" {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
```
