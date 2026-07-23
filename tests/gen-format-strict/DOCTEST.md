# Gen-format strict imports (Phase A) — L2 in-process

## Version
0.0.2

**Layer model:** all leaves are **L2 in-process** via `runner.RunTest` /
`core.WriteFormattedGo` (no product binary, no `testbin`). Leaves are
**unlabeled** (not `e2e`/`heavy`). Completeness: same six scenarios as before.

# DSN (Domain Specific Notion)

**Participants**

- **Author harness** — SETUP / Run / Assert Go blocks in a fixture doctest tree.
  Authors must **explicitly import** every package their own code names
  (`testing`, `time`, `session`, …). Omitted stdlib imports are **not** filled
  in by the generate write path.
- **Assemble** — lowers a fixture tree into hierarchical unified packages under
  a gen root (`__droot`, intermediate `setup.go`, leaf `RunTestLeaf`, suite).
  Only imports required **solely by doctest-injected harness** may be added
  (e.g. `session.Doctest` construct, `syscall` for locks, droot/registry).
- **WriteFormattedGo / format path** — post-assemble write for generated `.go`.
  Must **not** use `stdlibByPkgName` (or equivalent) to auto-add stdlib from
  selector usage, and must **not** require `go/format.Source` for a successful
  write that yields a **compilable** package.
- **Unused import prune** — assemble may over-import ancestor packages; the
  write path still drops imports that nothing in the file references.
- **Runner** — `runner.RunTest` / `build.Test` with explicit `GenDir` under
  `t.TempDir()` so generated sources are inspectable after generate (even when
  `go test` fails).

**Behaviors**

1. Fixture SETUP that uses `time.Second` with only `import "testing"` →
   generated package does **not** gain `"time"` → `go test` **fails**
   (`undefined: time` / undeclared name).
2. Same body with explicit `import "time"` → generate + suite `go test` **succeeds**.
3. Intermediate SETUP that does not reference a parent package’s symbols →
   generated intermediate `setup.go` does **not** retain the unused parent import.
4. Minimal leaf that only imports what author code uses (`testing`) still
   **compiles**: engine may add harness-only imports (`session`, etc.).
5. Synthetic source written via `WriteFormattedGo` that references `session.X`
   without an import does **not** gain a session import from the format path
   (no user-path auto-inject of session via stdlib-style maps).
6. Generated packages need only **compile**; tests assert compile success with
   explicit imports and **do not** require gofmt-pretty output.

```
fixture tree (author SETUP/ASSERT)
  -> Discover + Assemble (unified)
  -> WriteFormattedGo (no stdlib auto-add; unused prune OK; no gofmt requirement)
  -> go test suite package
  -> Assert: RunErr / generated .go text
```

## Decision Tree

```
gen-format-strict/
├── no-auto-stdlib/                         user stdlib must be explicit
│   ├── missing-time-fails/                 A1: time.Second, no import "time" → fail
│   └── explicit-time-ok/                   A2: import "time" → pass
├── prune-unused-ancestor/                  unused ancestor import dropped
│   └── intermediate-drops-unused-parent/   A3: mid setup.go lacks unused parent
├── harness-inject/                         engine may add harness-only imports
│   └── minimal-leaf-compiles/              A4: user imports testing only → compile
├── no-user-session-auto/                   format path does not auto-add session
│   └── session-without-import-not-auto-added/  A5: WriteFormattedGo leaves session out
└── compile-not-gofmt/                      format.Source not required for compile
    └── explicit-imports-compile/           A6: explicit imports → compile; no gofmt assert
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `no-auto-stdlib/missing-time-fails` | A1 — `time.Second` without `import "time"` → no auto-add; suite build fails |
| `no-auto-stdlib/explicit-time-ok` | A2 — explicit `import "time"` → generate + suite pass |
| `prune-unused-ancestor/intermediate-drops-unused-parent` | A3 — intermediate `setup.go` omits unused parent package import |
| `harness-inject/minimal-leaf-compiles` | A4 — user SETUP only imports `testing`; harness inject still compiles |
| `no-user-session-auto/session-without-import-not-auto-added` | A5 — `WriteFormattedGo` does not inject session for bare `session.` usage |
| `compile-not-gofmt/explicit-imports-compile` | A6 — explicit imports compile; Assert never requires gofmt |

## How to Run

```sh
doctest vet ./tests/gen-format-strict/
doctest test ./tests/gen-format-strict/
# L2 in-process (unlabeled); no --label heavy required
```

```go
import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

// Request selects the measured operation. Leaves build fixtures or source in Setup.
type Request struct {
	// Op:
	//   "run_fixture"        — write FixtureDir tree, runner.RunTest with GenDir
	//   "write_format_build" — WriteFormattedGo(Source) then optional go test/build
	Op string

	// FixtureKind selects which temp tree writeFixture builds when FixtureDir empty:
	//   "missing-time" | "explicit-time" | "unused-parent" | "harness-minimal" | "explicit-compile"
	FixtureKind string

	FixtureDir string // abs path to fixture tree (Setup may set)
	GenDir     string // abs gen root; empty → t.TempDir()

	// write_format_build
	Source    string // raw Go source (package clause included)
	OutGoName string // file name under GenDir; default "pkg.go"
	WantBuild bool   // if true, run `go build .` on written package
}

// Response captures generate/run outcomes and key generated texts.
type Response struct {
	RunErr     string
	Stdout     string
	Stderr     string
	FixtureDir string
	GenDir     string

	// Key generated sources (best-effort walk of GenDir).
	LeafGo              string // first leaf non-test .go with RunTestLeaf
	IntermediateSetupGo string // first …/mid/setup.go or intermediate setup.go under fixture mid
	AllGoSnippets       string // concat of selected .go bodies for substring asserts

	// write_format_build
	FormattedPath    string
	FormattedSource  string
	BuildErr         string
	BuildOutput      string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "run_fixture":
		return runFixture(t, req, resp)
	case "write_format_build":
		return runWriteFormatBuild(t, req, resp)
	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func runFixture(t *testing.T, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	dir := req.FixtureDir
	if dir == "" {
		dir = writeFixture(t, req.FixtureKind)
	}
	resp.FixtureDir = dir

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
		RemoveTemp: false,
	}
	err := runner.RunTest(dir, opts)
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()
	if err != nil {
		resp.RunErr = err.Error()
	}
	fillGenSources(t, resp)
	// Compile/test failures are measured outcomes, not harness errors.
	return resp, nil
}

func runWriteFormatBuild(t *testing.T, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	if req.Source == "" {
		return nil, fmt.Errorf("write_format_build requires Source")
	}
	genDir := req.GenDir
	if genDir == "" {
		genDir = t.TempDir()
	}
	resp.GenDir = genDir
	name := req.OutGoName
	if name == "" {
		name = "pkg.go"
	}
	outPath := filepath.Join(genDir, name)
	if err := core.WriteFormattedGo(outPath, req.Source); err != nil {
		// Write path may still leave a file; record and continue for inspection.
		resp.RunErr = err.Error()
	}
	resp.FormattedPath = outPath
	data, _ := os.ReadFile(outPath)
	resp.FormattedSource = string(data)
	if resp.FormattedSource == "" {
		// Fall back to attempted content if write failed before create.
		resp.FormattedSource = req.Source
	}

	if req.WantBuild {
		// Minimal module so `go build` can typecheck/compile the package.
		// (Avoid exec.Command("go","test") — banned anti-pattern in harnesses.)
		mod := "module genformatstrict\n\ngo 1.21\n"
		if err := os.WriteFile(filepath.Join(genDir, "go.mod"), []byte(mod), 0644); err != nil {
			return resp, err
		}
		cmd := exec.Command("go", "build", ".")
		cmd.Dir = genDir
		cmd.Env = append(os.Environ(), "GOWORK=off")
		out, err := cmd.CombinedOutput()
		resp.BuildOutput = string(out)
		if err != nil {
			resp.BuildErr = err.Error()
		}
	}
	return resp, nil
}

func fillGenSources(t *testing.T, resp *Response) {
	t.Helper()
	if resp.GenDir == "" {
		return
	}
	var snippets []string
	_ = filepath.Walk(resp.GenDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		rel, _ := filepath.Rel(resp.GenDir, path)
		relSlash := filepath.ToSlash(rel)
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		body := string(data)
		base := filepath.Base(path)

		// Intermediate packages are setup.go under non-leaf path segments.
		if base == "setup.go" {
			// Prefer mid/ segment when present (unused-parent fixture).
			if strings.Contains(relSlash, "/mid/") || strings.HasSuffix(relSlash, "/mid/setup.go") {
				resp.IntermediateSetupGo = body
			} else if resp.IntermediateSetupGo == "" && !strings.Contains(relSlash, "__") {
				resp.IntermediateSetupGo = body
			}
		}
		if strings.Contains(body, "func RunTestLeaf") && resp.LeafGo == "" {
			resp.LeafGo = body
		}
		// Keep a bounded concat for broad substring checks.
		if len(snippets) < 12 {
			snippets = append(snippets, "// file: "+relSlash+"\n"+body)
		}
		return nil
	})
	resp.AllGoSnippets = strings.Join(snippets, "\n\n")
}

// writeFixture builds an isolated fixture tree for the given kind.
func writeFixture(t *testing.T, kind string) string {
	t.Helper()
	root := t.TempDir()
	fence := string([]byte{'`', '`', '`'})

	switch kind {
	case "missing-time":
		writeMinimalRoot(t, root, fence, rootRunGo())
		writeLeaf(t, root, "leaf", fence,
			`import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = 20 * time.Second
	_ = req
	return nil
}`,
			leafAssertPass(),
		)
	case "explicit-time", "explicit-compile":
		writeMinimalRoot(t, root, fence, rootRunGo())
		writeLeaf(t, root, "leaf", fence,
			`import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	_ = 20 * time.Second
	_ = req
	return nil
}`,
			leafAssertPass(),
		)
	case "harness-minimal":
		// Author imports only testing; engine may still inject session/syscall for harness.
		writeMinimalRoot(t, root, fence, rootRunGo())
		writeLeaf(t, root, "leaf", fence,
			`import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}`,
			leafAssertPass(),
		)
	case "unused-parent":
		// root → feature (intermediate, exports FeatureHelper) → mid (no parent use) → leaf
		writeMinimalRoot(t, root, fence, rootRunGoWithHelper())
		writeGroupingSetup(t, root, "feature", fence,
			`import "testing"

// FeatureHelper is exported from the feature intermediate package.
func FeatureHelper() int { return 42 }

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}`)
		writeGroupingSetup(t, root, "feature/mid", fence,
			`import "testing"

func Setup(t *testing.T, req *Request) error {
	// Intentionally does NOT call FeatureHelper — parent import should be pruned.
	_ = t
	_ = req
	return nil
}`)
		writeLeaf(t, root, "feature/mid/leaf", fence,
			`import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}`,
			leafAssertPass(),
		)
	default:
		t.Fatalf("unknown FixtureKind %q", kind)
	}
	return root
}

func rootRunGo() string {
	return `import "testing"

type Request struct{}
type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = t
	_ = req
	return &Response{}, nil
}`
}

func rootRunGoWithHelper() string {
	return `import "testing"

type Request struct{}
type Response struct{}

func RootMarker() string { return "root" }

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = RootMarker()
	_ = t
	_ = req
	return &Response{}, nil
}`
}

func leafAssertPass() string {
	return `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	_ = req
	_ = resp
}`
}

func writeMinimalRoot(t *testing.T, root, fence, runGo string) {
	t.Helper()
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))
	// Root SETUP: real body (not stub) so validation accepts the fixture.
	setup := "# Scenario\n\n**Feature**: fixture root setup\n\n" +
		fence + "\nfixture root\n" + fence + "\n\n## Steps\n1. root setup\n\n" +
		fence + "go\n" + `import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}` + "\n" + fence + "\n"
	testtree.WriteFile(t, root, "SETUP.md", setup)
}

func writeGroupingSetup(t *testing.T, root, rel, fence, goCode string) {
	t.Helper()
	body := "# Scenario\n\n**Feature**: grouping " + rel + "\n\n" +
		fence + "\ngrouping " + rel + "\n" + fence + "\n\n## Steps\n1. grouping setup\n\n" +
		fence + "go\n" + goCode + "\n" + fence + "\n"
	testtree.WriteFile(t, root, filepath.Join(rel, "SETUP.md"), body)
}

func writeLeaf(t *testing.T, root, rel, fence, setupGo, assertGo string) {
	t.Helper()
	setup := "# Scenario\n\n**Feature**: leaf " + rel + "\n\n" +
		fence + "\nleaf " + rel + "\n" + fence + "\n\n## Steps\n1. leaf setup\n\n" +
		fence + "go\n" + setupGo + "\n" + fence + "\n"
	assert := "## Expected\n- passes\n\n" + fence + "go\n" + assertGo + "\n" + fence + "\n"
	testtree.WriteFile(t, root, filepath.Join(rel, "SETUP.md"), setup)
	testtree.WriteFile(t, root, filepath.Join(rel, "ASSERT.md"), assert)
}

// containsImportTime reports whether src has a time import path.
func containsImportTime(src string) bool {
	return strings.Contains(src, `"time"`)
}

// containsSessionImport reports whether src imports doctest/session.
func containsSessionImport(src string) bool {
	return strings.Contains(src, "github.com/xhd2015/doctest/session")
}

// looksGofmtEqual is true when gofmt leaves src unchanged (for A6 negative use).
func looksGofmtEqual(src string) bool {
	out, err := format.Source([]byte(src))
	if err != nil {
		return false
	}
	return string(out) == src
}
```
