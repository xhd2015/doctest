# Experiment: `--experiment-ref-instead-of-inline` — P0 plumbing + P1 ref generation

## Version
0.0.2

Classic-TDD specification spanning:

- **P0** (sealed): parse/plumb `--experiment-ref-instead-of-inline` into
  `core.Options.ExperimentRefInsteadOfInline` (default false).
- **P1** (this extension): when the flag is **on**, generate a **package DAG**
  for simple trees — root package owns `Request`/`Response`/`Run` (+ root
  helpers); thin leaf `_test.go` files **import** root instead of inlining
  ancestor source. Flag **off** keeps classic per-leaf inline assembly only.

These tests do **not** implement production code. Expect **RED** on new P1
leaves until the implementer lands the separate ref assemble path.

**Out of scope (P2+):** deep multi-level Setup edge cases, methods on Request
defined outside root (hard error may wait), full-suite perf, assert package
rewrites.

# DSN (Domain Specific Notion)

### Participants

- **Caller** — `doctest test` CLI (or package entry) that parses argv into options.
- **Parse layer** — `runner.ParseTestOptions(args)` filling `core.Options`.
- **Options** — `ExperimentRefInsteadOfInline bool` (CLI:
  `--experiment-ref-instead-of-inline`; default **false**).
- **Classic generator** — existing `AssembleTestSource`: one package per leaf;
  root types/`Run` inlined into each leaf `_test.go` (Run as local closure).
- **Ref generator (P1)** — separate assemble path when the flag is on: root
  package file(s) under the gen root hold `Run`/types/helpers once; intermediate
  SETUP dirs may become packages; leaf packages import ancestors; leaf tests are
  thin (chdir, inject `DOCTEST_*`, Setup chain, Run, Assert).
- **Gen root** — explicit `--gen-dir` / `Options.GenDir` under `t.TempDir()` in
  tests so layout is inspectable. Production should isolate ref mode from warm
  classic mapping-gen cache (e.g. mode marker or `mapping-gen-ref`).
- **Help** — documents the flag as experimental (`tests/help/test-options`).

### Behaviors

- **Default off** — omit flag → field false → classic inline only.
- **Opt-in** — flag → field true → ref package DAG for simple trees.
- **Run once** — distinctive root helper / `Run` body appears **once** on disk
  under gen root when flag on; twice (per leaf) when flag off (classic).
- **Thin leaves** — flag on: leaf `_test.go` does not redefine root types /
  marker helper; imports the root package path.
- **Both leaves pass** — simple 2-leaf fixture exits successfully under flag on.
- **Announce (optional)** — stderr may mention `experiment: ref-instead-of-inline`
  (or similar) when flag on.

### Pipeline sketch

```
doctest test [--experiment-ref-instead-of-inline] [--gen-dir DIR] <tree>
  -> ParseTestOptions -> Options.ExperimentRefInsteadOfInline
  -> if false: classic AssembleTestSource per leaf
  -> if true:  ref package DAG (root package + thin leaf tests)
  -> go test packages under gen root
```

## Decision Tree

```
tests/test/experiment-ref-inline/
├── flags/                                 [P0 ParseTestOptions]
│   ├── default-off/                       omit flag → false
│   └── flag-on/                           flag → true; remain has path
├── smoke/                                 [P0 mini RunTest]
│   ├── flag-off-still-passes/             field false → tiny tree ok
│   └── flag-on-still-passes/              field true → tiny tree ok
└── ref-mode/                              [P1 gen layout + 2-leaf fixture]
    ├── flag-off/                          classic path
    │   └── classic-inline-layout/         per-leaf inline; marker helper ×N leaves
    └── flag-on/                           ref path
        ├── two-leaves-pass/               exit/run success for a+b
        ├── root-marker-once/              marker helper defined once under gen
        ├── leaf-thin-import/              leaves import root; no type Request/marker def
        └── stderr-announce/               optional meta line on stderr
```

Sibling help (parent CLI tree): `tests/help/test-options`.

## Test Index

| Leaf | Phase | Expected |
|------|--------|----------|
| `flags/default-off` | P0 | field false by default |
| `flags/flag-on` | P0 | flag → true; remain includes path |
| `smoke/flag-off-still-passes` | P0 | mini run flag off ok |
| `smoke/flag-on-still-passes` | P0 | mini run flag on ok |
| `ref-mode/flag-off/classic-inline-layout` | P1 | 2 leaves pass; each leaf `_test.go` defines marker helper; no single shared root-only package required |
| `ref-mode/flag-on/two-leaves-pass` | P1 | 2-leaf fixture passes with flag on |
| `ref-mode/flag-on/root-marker-once` | P1 | `func ExperimentP1RootMarker` defined exactly once under gen |
| `ref-mode/flag-on/leaf-thin-import` | P1 | leaf `_test.go` lacks marker def + `type Request`; imports non-stdlib package |
| `ref-mode/flag-on/stderr-announce` | P1 | stderr contains `experiment` + `ref-instead-of-inline` (soft layout) |

## How to Run

```sh
doctest vet ./tests/test/experiment-ref-inline/
doctest test ./tests/test/experiment-ref-inline/
doctest test ./tests/test/experiment-ref-inline/flags/...
doctest test ./tests/test/experiment-ref-inline/smoke/...
doctest test ./tests/test/experiment-ref-inline/ref-mode/...
doctest test ./tests/help/test-options
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

// Distinctive root helper name/body used in P1 fixtures so gen layout can be counted.
const experimentP1MarkerFunc = "ExperimentP1RootMarker"
const experimentP1MarkerLiteral = "ROOT_RUN_MARKER_P1_EXPERIMENT_REF"

// Request selects one surface. Leaves set Op and related fields.
type Request struct {
	Op string // parse_flags | mini_run | ref_gen

	// parse_flags
	Args []string

	// mini_run / ref_gen
	ExperimentRefInsteadOfInline bool
	Dir                          string // fixture tree; empty → Run builds default fixture
	GenDir                       string // ref_gen only; empty → t.TempDir()
}

type Response struct {
	// flags
	Opts       core.Options
	RemainArgs []string
	ParseErr   string

	// mini_run / ref_gen
	RunErr   string
	Stdout   string
	Stderr   string
	Dir      string
	GenDir   string

	// ref_gen layout hints (filled by Run for assert convenience)
	GoFiles          []string // absolute paths of *.go under GenDir
	MarkerDefCount   int      // files defining func ExperimentP1RootMarker
	MarkerDefFiles   []string
	LeafTestFiles    []string // *a*_test.go / *b*_test.go style leaf tests
	LeafHasMarkerDef []bool   // parallel to LeafTestFiles
	LeafHasTypeReq   []bool
	LeafImportLines  [][]string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "parse_flags":
		opts, remain, err := runner.ParseTestOptions(req.Args)
		if err != nil {
			resp.ParseErr = err.Error()
			return resp, nil
		}
		resp.Opts = opts
		resp.RemainArgs = remain
		return resp, nil

	case "mini_run":
		dir := req.Dir
		if dir == "" {
			dir = createOnePassTree(t)
		}
		resp.Dir = dir

		var stdout, stderr bytes.Buffer
		opts := core.Options{
			Stdout:                       &stdout,
			Stderr:                       &stderr,
			RemoveTemp:                   true,
			ExperimentRefInsteadOfInline: req.ExperimentRefInsteadOfInline,
		}
		err := runner.RunTest(dir, opts)
		resp.Stdout = stdout.String()
		resp.Stderr = stderr.String()
		if err != nil {
			resp.RunErr = err.Error()
		}
		return resp, nil

	case "ref_gen":
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
			Stdout:                       &stdout,
			Stderr:                       &stderr,
			GenDir:                       genDir,
			RemoveTemp:                   false, // keep GenDir for layout asserts
			ExperimentRefInsteadOfInline: req.ExperimentRefInsteadOfInline,
		}
		err := runner.RunTest(dir, opts)
		resp.Stdout = stdout.String()
		resp.Stderr = stderr.String()
		if err != nil {
			resp.RunErr = err.Error()
		}
		fillGenLayout(t, resp)
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func createOnePassTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 1, 0)
	return tmp
}

// createTwoLeafMarkerTree builds a simple 2-leaf doctest tree whose root DOCTEST
// defines ExperimentP1RootMarker + Run that references it. Classic gen duplicates
// the helper per leaf; ref gen should emit it once in a shared root package.
func createTwoLeafMarkerTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	runGo := `import "testing"

type Request struct{}
type Response struct{}

// ExperimentP1RootMarker is a distinctive root helper for gen-layout asserts.
func ExperimentP1RootMarker() string {
	return "ROOT_RUN_MARKER_P1_EXPERIMENT_REF"
}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = ExperimentP1RootMarker()
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

func fillGenLayout(t *testing.T, resp *Response) {
	t.Helper()
	if resp.GenDir == "" {
		return
	}
	var goFiles []string
	_ = filepath.Walk(resp.GenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	resp.GoFiles = goFiles

	markerSig := "func " + experimentP1MarkerFunc
	for _, path := range goFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		body := string(data)
		if strings.Contains(body, markerSig) {
			resp.MarkerDefCount++
			resp.MarkerDefFiles = append(resp.MarkerDefFiles, path)
		}
		base := filepath.Base(path)
		// Leaf tests for fixture leaves a/ and b/ (classic: a/a_test.go, b/b_test.go).
		if strings.HasSuffix(base, "_test.go") {
			rel, _ := filepath.Rel(resp.GenDir, path)
			relSlash := filepath.ToSlash(rel)
			if strings.Contains(relSlash, "/a/") || strings.HasPrefix(relSlash, "a/") ||
				strings.Contains(relSlash, "/b/") || strings.HasPrefix(relSlash, "b/") ||
				base == "a_test.go" || base == "b_test.go" {
				resp.LeafTestFiles = append(resp.LeafTestFiles, path)
				resp.LeafHasMarkerDef = append(resp.LeafHasMarkerDef, strings.Contains(body, markerSig))
				resp.LeafHasTypeReq = append(resp.LeafHasTypeReq, strings.Contains(body, "type Request"))
				resp.LeafImportLines = append(resp.LeafImportLines, importLines(body))
			}
		}
	}
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

func leafHasNonStdImport(importLines []string) bool {
	for _, line := range importLines {
		// strip alias
		s := strings.Trim(line, "\t ")
		if s == "" || strings.HasPrefix(s, "//") {
			continue
		}
		// `"path"` or `alias "path"`
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
			// any other path counts (root package under testcase module)
			if path != "" {
				return true
			}
		}
	}
	return false
}
```
