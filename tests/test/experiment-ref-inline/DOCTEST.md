# Default hierarchical ref packages (under unified suite generation)

## Version
0.0.3

Default generation uses a **hierarchical ref package DAG** inside the unified
suite path: root package owns `Request`/`Response`/`Run` (+ root helpers);
leaves import ancestors instead of inlining root source.

No experiment flags; classic full-inline per leaf is not the production default.

# DSN (Domain Specific Notion)

### Participants

- **Caller** — `doctest test` CLI (or package entry).
- **Ref/unified generator** — root package under `__droot` holds `Run`/types/helpers
  once; intermediate SETUP dirs may become packages; leaf packages import
  ancestors; unified leaves expose `RunTestLeaf` (non-test packages).
- **Gen root** — explicit `Options.GenDir` under `t.TempDir()` so layout is
  inspectable. Production warm cache is `mapping-gen`.

### Behaviors

- **Default** — hierarchical ref packages (shared root marker once).
- **Run once** — distinctive root helper / `Run` body appears **once** on disk
  under gen root.
- **Thin leaves** — leaf packages do not redefine root types / marker helper;
  import non-stdlib ancestor packages.
- **Both leaves pass** — simple 2-leaf fixture exits successfully under default.
- **Smoke** — mini one-leaf tree still passes under default generation.

### Pipeline sketch

```
doctest test [--gen-dir DIR] <tree>
  -> hierarchical unified (ref packages + suite)
  -> go test suite package under gen root
```

## Decision Tree

```
tests/test/experiment-ref-inline/
├── smoke/                                 [mini RunTest]
│   └── default-passes/                    tiny tree ok under default gen
└── ref-mode/                              [gen layout + 2-leaf fixture]
    ├── two-leaves-pass/                   exit/run success for a+b
    ├── root-marker-once/                  marker helper defined once under gen
    └── leaf-thin-import/                  leaves import root; no type Request/marker def
```

## Test Index

| Leaf | Expected |
|------|----------|
| `smoke/default-passes` | mini run under default gen ok |
| `ref-mode/two-leaves-pass` | 2-leaf fixture passes under default |
| `ref-mode/root-marker-once` | `func ExperimentP1RootMarker` defined exactly once under gen |
| `ref-mode/leaf-thin-import` | leaf packages lack marker def + `type Request`; import non-stdlib package |

## How to Run

```sh
doctest vet ./tests/test/experiment-ref-inline/
doctest test ./tests/test/experiment-ref-inline/
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

// Distinctive root helper name/body used in fixtures so gen layout can be counted.
const experimentP1MarkerFunc = "ExperimentP1RootMarker"
const experimentP1MarkerLiteral = "ROOT_RUN_MARKER_P1_EXPERIMENT_REF"

// Request selects one surface. Leaves set Op and related fields.
type Request struct {
	Op string // mini_run | ref_gen

	// mini_run / ref_gen
	Dir    string // fixture tree; empty → Run builds default fixture
	GenDir string // ref_gen only; empty → t.TempDir()
}

type Response struct {
	// mini_run / ref_gen
	RunErr string
	Stdout string
	Stderr string
	Dir    string
	GenDir string

	// ref_gen layout hints (filled by Run for assert convenience)
	GoFiles          []string // absolute paths of *.go under GenDir
	MarkerDefCount   int      // files defining func ExperimentP1RootMarker
	MarkerDefFiles   []string
	LeafGoFiles      []string // leaf package .go under a/ or b/ (non-suite)
	LeafHasMarkerDef []bool   // parallel to LeafGoFiles
	LeafHasTypeReq   []bool
	LeafImportLines  [][]string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "mini_run":
		dir := req.Dir
		if dir == "" {
			dir = createOnePassTree(t)
		}
		resp.Dir = dir

		var stdout, stderr bytes.Buffer
		opts := core.Options{
			Stdout:     &stdout,
			Stderr:     &stderr,
			RemoveTemp: true,
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
// defines ExperimentP1RootMarker + Run that references it. Hierarchical gen should
// emit the helper once in the shared root package.
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
		rel, _ := filepath.Rel(resp.GenDir, path)
		relSlash := filepath.ToSlash(rel)
		// Fixture leaves a/ and b/ only (not suite/registry/allleaves/droot).
		if !isLeafABPath(relSlash) {
			continue
		}
		resp.LeafGoFiles = append(resp.LeafGoFiles, path)
		resp.LeafHasMarkerDef = append(resp.LeafHasMarkerDef, strings.Contains(body, markerSig))
		resp.LeafHasTypeReq = append(resp.LeafHasTypeReq, strings.Contains(body, "type Request"))
		resp.LeafImportLines = append(resp.LeafImportLines, importLines(body))
	}
}

func isLeafABPath(relSlash string) bool {
	parts := strings.Split(relSlash, "/")
	for i, p := range parts {
		if p == "a" || p == "b" {
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
