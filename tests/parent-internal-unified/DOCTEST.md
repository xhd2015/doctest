# Parent-internal always unified (layout A) — P2

## Version

0.0.2

## Layer

**L2 in-process** — `runner.RunTest` + `core.Options` on a temp parent module
with multi-leaf tree importing `example.com/app/internal/greet`. No product
binary. Classic TDD for P2: leaves are **RED** while parent-internal still
forces `internalCompile` multi-leaf under `.doctest_run_*`; **GREEN** after
always-unified mapping-gen (suite package + Kind B expose).

# DSN (Domain Specific Notion)

### Participants

- **Caller** — this harness via `runner.RunTest` / library prepare+test path.
- **Parent module** — temp `example.com/app` with exported `internal/greet`
  (`Hello` + `DefaultName` var).
- **Subject tree** — multi-leaf doctest under `tests/` importing product internal.
- **Unified generator** — hierarchical mapping-gen (layout A): `__droot`,
  `__registry`, leaf `RunTestLeaf`, `__allleaves`, single `suite` package.
- **Kind B expose** — overlay facades + import rewrite so unified gen can load
  product internal from module `testcase` (implementation detail).

### Behaviors

- **Always unified** — parent-internal trees do **not** default to multi-leaf
  `internalCompile` / `.doctest_run_*` package lists.
- **Subject pass** — both leaves exercise `Hello` and `DefaultName` and pass.
- **Suite-only go test** — displayed package args are a single path containing
  `suite`, not `./leaf-a ./leaf-b`.
- **Gen layout** — gen dir has `suite` + `__allleaves` (unified markers).
- **Coverprofile** — `-cover` + `-coverprofile` exit 0; profile file non-empty
  (single package).

### Pipeline sketch

```
temp module example.com/app + internal/greet + multi-leaf tests/
  -> runner.RunTest(tests, GenDir, optional CoverProfile)
  -> unified mapping-gen + Kind B expose
  -> go test ./…/suite  (one package)
  -> both leaves PASS; cover.out written when requested
```

## Decision Tree

```
parent-internal-unified/
└── multi-leaf/                         in-module ≥2 leaves → parent internal
    ├── subject-and-suite/              both leaves PASS + suite-only go test
    ├── gen-layout/                     suite + __allleaves (unified shape)
    └── coverprofile/                   -cover -coverprofile exit 0 + file
```

## Test Index

| # | Leaf | Expected after P2 (RED today) |
|---|------|-------------------------------|
| 1 | `multi-leaf/subject-and-suite` | Subject tests PASS; go test single suite package (not multi leaf) |
| 2 | `multi-leaf/gen-layout` | Gen has `suite` + `__allleaves`; subject run succeeds |
| 3 | `multi-leaf/coverprofile` | Coverprofile succeeds; file exists non-empty |

## How to Run

```sh
doctest vet ./tests/parent-internal-unified
doctest test ./tests/parent-internal-unified
doctest test -v ./tests/parent-internal-unified/multi-leaf/coverprofile
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
	"github.com/xhd2015/doctest/session"
)

const modPath = "example.com/app"

// Request selects the acceptance facet for one multi-leaf parent-internal run.
type Request struct {
	// WithCover enables -cover + CoverProfile under req.CoverPath.
	WithCover bool
	// CoverPath is absolute path for -coverprofile (filled by Setup when WithCover).
	CoverPath string
	// ModuleRoot / TestDir / GenDir filled by multi-leaf Setup.
	ModuleRoot string
	TestDir    string
	GenDir     string
}

type Response struct {
	RunErr string
	Stdout string
	Stderr string

	ModuleRoot string
	TestDir    string
	GenDir     string
	CoverPath  string

	// Layout / go-test packaging (filled after RunTest).
	HasSuite     bool
	HasAllLeaves bool
	HasDroot     bool
	HasRegistry  bool
	GoFiles      []string

	GoTestDisplayLine string
	GoTestPackageArgs []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req.TestDir == "" || req.GenDir == "" {
		return nil, fmt.Errorf("TestDir and GenDir must be set by multi-leaf Setup")
	}
	resp := &Response{
		ModuleRoot: req.ModuleRoot,
		TestDir:    req.TestDir,
		GenDir:     req.GenDir,
		CoverPath:  req.CoverPath,
	}

	var stdout, stderr bytes.Buffer
	opts := core.Options{
		Stdout:     &stdout,
		Stderr:     &stderr,
		GenDir:     req.GenDir,
		RemoveTemp: false, // keep GenDir for layout asserts
		Count:      1,
		Verbose:    true, // ensure go test line is easy to parse
	}
	if req.WithCover {
		if req.CoverPath == "" {
			return nil, fmt.Errorf("WithCover requires CoverPath")
		}
		opts.Cover = true
		opts.CoverProfile = req.CoverPath
	}

	err := runner.RunTest(req.TestDir, opts)
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()
	if err != nil {
		resp.RunErr = err.Error()
	}
	fillParentInternalLayout(t, resp)
	return resp, nil
}

func fillParentInternalLayout(t *testing.T, resp *Response) {
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
			if strings.Contains(relSlash, "__droot") {
				resp.HasDroot = true
			}
			if strings.Contains(relSlash, "__registry") {
				resp.HasRegistry = true
			}
			if strings.Contains(relSlash, "__allleaves") {
				resp.HasAllLeaves = true
			}
			if relSlash == "suite" || strings.HasPrefix(relSlash, "suite/") || strings.Contains(relSlash, "/suite/") {
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
	resp.GoTestDisplayLine, resp.GoTestPackageArgs = parseGoTestPackageArgs(resp.Stdout + "\n" + resp.Stderr)
}

// parseGoTestPackageArgs finds the first displayed `go test …` line and returns
// package path arguments (tokens starting with ./ or exactly ".").
func parseGoTestPackageArgs(combined string) (displayLine string, pkgs []string) {
	for _, line := range strings.Split(combined, "\n") {
		trim := strings.TrimSpace(line)
		idx := strings.Index(trim, "go test")
		if idx < 0 {
			continue
		}
		rest := strings.TrimSpace(trim[idx:])
		displayLine = rest
		fields := strings.Fields(rest)
		for i := 2; i < len(fields); i++ {
			f := fields[i]
			if strings.HasPrefix(f, "-") {
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

// createParentInternalMultiLeafModule builds:
//
//	module example.com/app
//	internal/greet (Hello + DefaultName)
//	tests/ leaf-a + leaf-b importing product internal
func createParentInternalMultiLeafModule(t *testing.T) (moduleRoot, testDir string) {
	t.Helper()
	moduleRoot = t.TempDir()
	if err := os.WriteFile(
		filepath.Join(moduleRoot, "go.mod"),
		[]byte("module "+modPath+"\n\ngo 1.21\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	greetDir := filepath.Join(moduleRoot, "internal", "greet")
	if err := os.MkdirAll(greetDir, 0755); err != nil {
		t.Fatalf("mkdir internal/greet: %v", err)
	}
	greetSrc := `package greet

// DefaultName is an exported var used by Hello when name is empty.
var DefaultName = "world"

// Hello returns a greeting; empty name uses DefaultName.
func Hello(name string) string {
	if name == "" {
		name = DefaultName
	}
	return "hello " + name
}
`
	if err := os.WriteFile(filepath.Join(greetDir, "greet.go"), []byte(greetSrc), 0644); err != nil {
		t.Fatalf("write greet.go: %v", err)
	}

	testDir = filepath.Join(moduleRoot, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}

	runGo := `import (
	"testing"

	"` + modPath + `/internal/greet"
	"github.com/xhd2015/doctest/session"
)

type Request struct {
	// Name is passed to greet.Hello. Empty selects DefaultName path.
	Name string
}

type Response struct {
	Message string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	return &Response{Message: greet.Hello(req.Name)}, nil
}`
	testtree.WriteFile(t, testDir, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))

	fence := string([]byte{'`', '`', '`'})

	// leaf-a: Hello("alice") exercises the function API.
	writeSubjectLeaf(t, testDir, "leaf-a", fence,
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Name = "alice"
	return nil
}`,
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Message != "hello alice" {
		t.Fatalf("Message = %v, want %q", resp, "hello alice")
	}
}`,
	)

	// leaf-b: Hello("") uses DefaultName ("world").
	writeSubjectLeaf(t, testDir, "leaf-b", fence,
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Name = ""
	return nil
}`,
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Message != "hello world" {
		t.Fatalf("Message = %v, want %q (via DefaultName)", resp, "hello world")
	}
}`,
	)

	return moduleRoot, testDir
}

func writeSubjectLeaf(t *testing.T, testDir, name, fence, setupGo, assertGo string) {
	t.Helper()
	setupMD := "# Scenario\n\n**Feature**: subject leaf " + name + "\n\n" +
		fence + "\nleaf " + name + " imports parent internal/greet\n" + fence + "\n\n## Steps\n1. leaf setup\n\n" +
		fence + "go\n" + setupGo + "\n" + fence + "\n"
	assertMD := "## Expected\n- subject leaf " + name + " passes\n\n" +
		fence + "go\n" + assertGo + "\n" + fence + "\n"
	testtree.WriteFile(t, testDir, name+"/SETUP.md", setupMD)
	testtree.WriteFile(t, testDir, name+"/ASSERT.md", assertMD)
}
```
