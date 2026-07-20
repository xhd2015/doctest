# Default: hierarchical unified package per DOCTEST tree

## Version
0.0.4

Default generation produces **one go test package/binary per DOCTEST tree**
using hierarchical ref packages (no experiment flags).

Per DOCTEST.md root, generate:

1. `__droot/` — shared types/Run (ref)
2. `__registry/` — tree-local `Register` / `All`
3. leaf packages as **non-`_test`** sources with **`RunTestLeaf(t *testing.T)`** + `init()` registration
4. `__allleaves/` — blank-imports every leaf package
5. `suite/` — package with importable `RunAll` (`runall.go`: registry + blank
   `__allleaves`) and thin `suite_test.go` (`TestDoctestSuite` → `RunAll`; stdlib only)

**`go test` only the suite package** → **one test binary** per tree.
Suite **package** non-stdlib imports are only `__registry` + `__allleaves`
(fan-in lives in `runall.go`, not necessarily in `*_test.go`).

# DSN (Domain Specific Notion)

### Participants

- **Caller** — `doctest test` CLI (or package entry).
- **Unified generator** — hierarchical ref + leaf **non-test** packages with
  `RunTestLeaf`; `__registry` + `__allleaves` fan-in; single `suite` package
  with `RunAll` + thin `TestDoctestSuite` (hub-importable).
- **Gen root** — explicit `--gen-dir` / `Options.GenDir` under `t.TempDir()` so
  layout is inspectable.

### Behaviors

- **Default** — no flag → unified hierarchical layout.
- **Layout** — gen contains `__droot`, `__registry`, `__allleaves`, `suite`;
  leaves are non-`_test` with `RunTestLeaf`; no per-leaf `*_test.go` for fixture
  leaves `a`/`b`.
- **Single package** — displayed `go test` line lists only the suite package
  (not `./a ./b`).
- **Both leaves pass** — simple 2-leaf fixture exits successfully under default.

### Pipeline sketch

```
doctest test [--gen-dir DIR] <tree>
  -> hierarchical unified gen: __droot + __registry + leaf RunTestLeaf + __allleaves + suite (runall + suite_test)
  -> go test ./…/suite   (one package / one binary)
```

## Decision Tree

```
tests/test/experiment-unified-package/
└── unified-mode/                          [RunTest default + GenDir]
    ├── two-leaves-pass/                   exit/run success for a+b
    ├── gen-layout/                        __droot/__registry/__allleaves/suite;
    │                                      leaf RunTestLeaf; no leaf *_test.go
    └── suite-only-go-test/                go test line: single suite package
```

## Test Index

| Leaf | Expected |
|------|----------|
| `unified-mode/two-leaves-pass` | 2-leaf fixture passes under default gen |
| `unified-mode/gen-layout` | dirs + leaf `RunTestLeaf`; suite package imports only registry/allleaves (`runall` + thin test) |
| `unified-mode/suite-only-go-test` | displayed `go test` has one package containing `suite` (not `./a ./b`) |

## How to Run

```sh
doctest vet ./tests/test/experiment-unified-package/
doctest test ./tests/test/experiment-unified-package/
```

```go
import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Distinctive root helper so layout can be counted without hard-coding paths.
const experimentUnifiedMarkerFunc = "ExperimentUnifiedRootMarker"
const experimentUnifiedMarkerLiteral = "ROOT_RUN_MARKER_UNIFIED_PACKAGE"

// Request selects one surface. Leaves set Op and related fields.
type Request struct {
	Op string // run_gen

	// run_gen
	Dir    string // fixture tree; empty → Run builds default
	GenDir string // empty → t.TempDir()
}

type Response struct {
	// run_gen
	RunErr string
	Stdout string
	Stderr string
	Dir    string
	GenDir string

	// layout / go-test package hints (filled by Run)
	GoFiles            []string // absolute paths of *.go under GenDir
	HasDroot           bool
	HasRegistry        bool
	HasAllLeaves       bool
	HasSuite           bool
	SuiteTestFiles     []string // *_test.go under …/suite/
	SuiteNonTestFiles  []string // non-_test .go under …/suite/ (e.g. runall.go)
	SuiteImportLines   []string // imports from all suite package .go files
	LeafNonTestGoFiles []string // non-_test .go under leaf a/ or b/
	LeafHasRunTestLeaf []bool   // parallel to LeafNonTestGoFiles
	LeafTestGoFiles    []string // *_test.go under leaf a/ or b/ (want empty when unified)
	GoTestPackageArgs  []string // ./… package args from displayed go test line
	GoTestDisplayLine  string   // full "go test …" line from stdout/stderr
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "run_gen":
		dir := req.Dir
		if dir == "" {
			dir = createTwoLeafMarkerTree(t)
		}
		resp.Dir = dir

		genDir := req.GenDir
		if genDir == "" {
			genDir = t.TempDir()
		}
		resp.GenDir = genDir

		var stdout, stderr bytes.Buffer
		opts := core.Options{
			Stdout:     &stdout,
			Stderr:     &stderr,
			GenDir:     genDir,
			RemoveTemp: false, // keep GenDir for layout asserts
		}
		err := runner.RunTest(dir, opts)
		resp.Stdout = stdout.String()
		resp.Stderr = stderr.String()
		if err != nil {
			resp.RunErr = err.Error()
		}
		fillUnifiedGenLayout(t, resp)
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

// createTwoLeafMarkerTree builds a simple 2-leaf doctest tree whose root DOCTEST
// defines ExperimentUnifiedRootMarker + Run that references it.
func createTwoLeafMarkerTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	runGo := `import "testing"

type Request struct{}
type Response struct{}

// ExperimentUnifiedRootMarker is a distinctive root helper for gen-layout asserts.
func ExperimentUnifiedRootMarker() string {
	return "ROOT_RUN_MARKER_UNIFIED_PACKAGE"
}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = ExperimentUnifiedRootMarker()
	return &Response{}, nil
}`
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))

	// Build md fences without embedding triple-backtick sequences in this
	// DOCTEST.md Go block (would confuse ExtractFinalGoBlock).
	fence := string([]byte{'`', '`', '`'})
	for _, name := range []string{"a", "b"} {
		setup := `import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}`
		assert := `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
}`
		setupMD := "# Scenario\n\n**Feature**: leaf " + name + "\n\n" +
			fence + "\nleaf " + name + "\n" + fence + "\n\n## Steps\n1. leaf setup\n\n" +
			fence + "go\n" + setup + "\n" + fence + "\n"
		assertMD := "## Expected\n- passes\n\n" + fence + "go\n" + assert + "\n" + fence + "\n"
		testtree.WriteFile(t, root, name+"/SETUP.md", setupMD)
		testtree.WriteFile(t, root, name+"/ASSERT.md", assertMD)
	}
	return root
}

func fillUnifiedGenLayout(t *testing.T, resp *Response) {
	t.Helper()
	if resp.GenDir == "" {
		return
	}
	var goFiles []string
	_ = filepath.Walk(resp.GenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		rel, relErr := filepath.Rel(resp.GenDir, path)
		if relErr != nil {
			rel = path
		}
		relSlash := filepath.ToSlash(rel)
		if info.IsDir() {
			base := filepath.Base(path)
			switch base {
			case "__droot":
				resp.HasDroot = true
			case "__registry":
				resp.HasRegistry = true
			case "__allleaves":
				resp.HasAllLeaves = true
			case "suite":
				resp.HasSuite = true
			}
			// also match path segments (tree-scoped layouts)
			if strings.Contains(relSlash, "__droot") {
				resp.HasDroot = true
			}
			if strings.Contains(relSlash, "__registry") {
				resp.HasRegistry = true
			}
			if strings.Contains(relSlash, "__allleaves") {
				resp.HasAllLeaves = true
			}
			if strings.Contains(relSlash, "/suite") || relSlash == "suite" || strings.HasPrefix(relSlash, "suite/") {
				resp.HasSuite = true
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	resp.GoFiles = goFiles

	for _, path := range goFiles {
		rel, _ := filepath.Rel(resp.GenDir, path)
		relSlash := filepath.ToSlash(rel)
		base := filepath.Base(path)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		body := string(data)

		isSuite := strings.Contains(relSlash, "/suite/") || strings.HasPrefix(relSlash, "suite/") ||
			base == "suite_test.go" || base == "runall.go"
		if isSuite && strings.HasSuffix(base, ".go") {
			resp.HasSuite = true
			if strings.HasSuffix(base, "_test.go") {
				resp.SuiteTestFiles = append(resp.SuiteTestFiles, path)
			} else {
				resp.SuiteNonTestFiles = append(resp.SuiteNonTestFiles, path)
			}
			// Package-level: merge imports from runall.go + suite_test.go (Option C /
			// hub layout — fan-in lives in RunAll, not only in *_test.go).
			resp.SuiteImportLines = append(resp.SuiteImportLines, importLines(body)...)
		}

		// Fixture leaves a/ and b/ only (not suite/registry/allleaves/droot).
		if isLeafABPath(relSlash) {
			if strings.HasSuffix(base, "_test.go") {
				resp.LeafTestGoFiles = append(resp.LeafTestGoFiles, path)
			} else if strings.HasSuffix(base, ".go") {
				resp.LeafNonTestGoFiles = append(resp.LeafNonTestGoFiles, path)
				resp.LeafHasRunTestLeaf = append(resp.LeafHasRunTestLeaf,
					strings.Contains(body, "func RunTestLeaf"))
			}
		}
	}

	resp.GoTestDisplayLine, resp.GoTestPackageArgs = parseGoTestPackageArgs(resp.Stdout + "\n" + resp.Stderr)
}

func isLeafABPath(relSlash string) bool {
	// Match gen paths like a/foo.go, b/b.go, tree/a/x.go — exclude suite etc.
	parts := strings.Split(relSlash, "/")
	for i, p := range parts {
		if p == "a" || p == "b" {
			// must not be under __* or suite special packages
			for _, prev := range parts[:i] {
				if prev == "suite" || strings.HasPrefix(prev, "__") {
					return false
				}
			}
			return true
		}
	}
	return false
}

func importLines(src string) []string {
	var lines []string
	inImport := false
	for _, line := range strings.Split(src, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "import (") {
			inImport = true
			continue
		}
		if inImport {
			if trim == ")" {
				inImport = false
				continue
			}
			lines = append(lines, trim)
			continue
		}
		if strings.HasPrefix(trim, "import ") {
			lines = append(lines, trim)
		}
	}
	return lines
}

// parseGoTestPackageArgs finds the first displayed `go test …` line and returns
// package path arguments (tokens starting with ./ or exactly ".").
func parseGoTestPackageArgs(combined string) (displayLine string, pkgs []string) {
	for _, line := range strings.Split(combined, "\n") {
		trim := strings.TrimSpace(line)
		// Match "cd X && go test …" or bare "go test …"
		idx := strings.Index(trim, "go test")
		if idx < 0 {
			continue
		}
		rest := strings.TrimSpace(trim[idx:])
		displayLine = rest
		fields := strings.Fields(rest)
		// skip "go" "test"
		for i := 2; i < len(fields); i++ {
			f := fields[i]
			if strings.HasPrefix(f, "-") {
				// flag; may be -flag=value already one field
				continue
			}
			if f == "." || strings.HasPrefix(f, "./") {
				pkgs = append(pkgs, f)
			}
		}
		return displayLine, pkgs
	}
	return "", nil
}

func suiteImportsOnlyRegistryAndAllLeaves(importLines []string) bool {
	var nonStd []string
	for _, line := range importLines {
		s := strings.Trim(line, "\t ")
		if s == "" || strings.HasPrefix(s, "//") {
			continue
		}
		q := strings.Index(s, `"`)
		if q < 0 {
			continue
		}
		rest := s[q+1:]
		end := strings.Index(rest, `"`)
		if end < 0 {
			continue
		}
		path := rest[:end]
		switch path {
		case "testing", "os", "path/filepath", "syscall", "fmt", "strings", "bytes", "time", "context", "errors":
			continue
		default:
			if path != "" {
				nonStd = append(nonStd, path)
			}
		}
	}
	if len(nonStd) == 0 {
		return false
	}
	sawRegistry, sawAllLeaves := false, false
	for _, p := range nonStd {
		if strings.Contains(p, "__registry") {
			sawRegistry = true
			continue
		}
		if strings.Contains(p, "__allleaves") {
			sawAllLeaves = true
			continue
		}
		// any other non-stdlib import is not allowed for suite
		return false
	}
	return sawRegistry && sawAllLeaves
}
```
