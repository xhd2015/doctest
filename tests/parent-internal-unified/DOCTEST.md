# Parent-internal always unified (layout A) — P2

## Version

0.0.2

## Layer

**L2 in-process** — `runner.RunTest` + `core.Options` on a temp parent module
with multi-leaf tree importing `example.com/app/internal/greet`. No product
binary. Classic TDD for P2: leaves are **RED** while parent-internal still
forces `internalCompile` multi-leaf under `.doctest_run_*`; **GREEN** after
always-unified mapping-gen (suite package + expose).

# DSN (Domain Specific Notion)

### Participants

- **Caller** — this harness via `runner.RunTest` / library prepare+test path.
- **Parent module** — temp `example.com/app` with exported `internal/greet`
  (`Hello` + `DefaultName` var).
- **Subject tree** — multi-leaf doctest under `tests/` importing product internal.
- **Unified generator** — hierarchical mapping-gen (layout A): `__droot`,
  `__registry`, leaf `RunTestLeaf`, `__allleaves`, single `suite` package.
- **Expose** — overlay facades + import rewrite so unified gen can load
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
- **Coverpkg + expose** — with `-coverpkg=example.com/app/...` (product
  module wildcard; same shape as CI `github.com/<mod>/...`), cover instruments
  product packages; suite run must not fail opening
  `…/__doctest_internal_expose/…/expose.go`. The **final** `-coverprofile`
  must omit session-generated expose facade paths so
  `go tool cover -func=<profile>` from the product module succeeds (scaff CI
  merge/report).
- **External types in internal signatures** — when product `internal/…` exports
  funcs using types from another product package (e.g. `model.Project`),
  expose facades must compile (import or re-alias those packages). No
  `undefined: model` (or twin) on the generated expose body.

### Pipeline sketch

```
temp module example.com/app + internal/greet + multi-leaf tests/
  -> runner.RunTest(tests, GenDir, optional CoverProfile[, CoverPkg])
  -> unified mapping-gen + expose
  -> go test ./…/suite  (one package)
  -> both leaves PASS; cover.out written when requested
```

## Decision Tree

```
parent-internal-unified/
├── multi-leaf/                         in-module ≥2 leaves → parent internal
│   ├── subject-and-suite/              both leaves PASS + suite-only go test
│   ├── gen-layout/                     suite + __allleaves (unified shape)
│   ├── coverprofile/                   -cover -coverprofile exit 0 + file
│   └── coverpkg-expose/                -coverpkg=mod/...; clean profile + go tool cover
└── external-sig-types/                 internal API uses model.T; expose must compile
```

## Test Index

| # | Leaf | Expected after P2 (RED today) |
|---|------|-------------------------------|
| 1 | `multi-leaf/subject-and-suite` | Subject tests PASS; go test single suite package (not multi leaf) |
| 2 | `multi-leaf/gen-layout` | Gen has `suite` + `__allleaves`; subject run succeeds |
| 3 | `multi-leaf/coverprofile` | Coverprofile succeeds; file exists non-empty |
| 4 | `multi-leaf/coverpkg-expose` | Cover + coverpkg succeeds; profile omits expose facades; `go tool cover -func` OK |
| 5 | `external-sig-types` | expose compiles when internal uses external package types |

## How to Run

```sh
doctest vet ./tests/parent-internal-unified
doctest test ./tests/parent-internal-unified
doctest test -v ./tests/parent-internal-unified/multi-leaf/coverprofile
doctest test -v ./tests/parent-internal-unified/multi-leaf/coverpkg-expose
doctest test -v ./tests/parent-internal-unified/external-sig-types
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
	// CoverPkg is go test -coverpkg (comma-separated). Empty = omit.
	// Use product module wildcards (e.g. example.com/app/...) to match CI.
	CoverPkg string
	// CoverMode is go test -covermode (set|count|atomic). Empty = omit.
	CoverMode string
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
	if req.CoverPkg != "" {
		opts.CoverPkg = req.CoverPkg
		// coverpkg alone implies coverage analysis; keep Cover true for clarity.
		opts.Cover = true
	}
	if req.CoverMode != "" {
		opts.CoverMode = req.CoverMode
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

// createParentInternalExternalSigModule builds a parent module where
// internal/rules exported API uses types from product package model
// (crime scene: expose must import model or re-alias types).
//
//	module example.com/app
//	model (Project, FixResult)
//	internal/rules.FixIgnore(model.Project, bool) (model.FixResult, error)
//	tests/leaf-a subject importing internal/rules
func createParentInternalExternalSigModule(t *testing.T) (moduleRoot, testDir string) {
	t.Helper()
	moduleRoot = t.TempDir()
	if err := os.WriteFile(
		filepath.Join(moduleRoot, "go.mod"),
		[]byte("module "+modPath+"\n\ngo 1.21\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	modelDir := filepath.Join(moduleRoot, "model")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		t.Fatalf("mkdir model: %v", err)
	}
	modelSrc := `package model

// Project is an external product type (stand-in for scaff model.Project).
type Project struct {
	Root string
}

// FixResult is returned from fix helpers.
type FixResult struct {
	OK bool
}
`
	if err := os.WriteFile(filepath.Join(modelDir, "project.go"), []byte(modelSrc), 0644); err != nil {
		t.Fatalf("write model: %v", err)
	}

	rulesDir := filepath.Join(moduleRoot, "internal", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		t.Fatalf("mkdir internal/rules: %v", err)
	}
	rulesSrc := `package rules

import "` + modPath + `/model"

// FixIgnore uses external package types in the exported signature.
func FixIgnore(project model.Project, dryRun bool) (model.FixResult, error) {
	_ = dryRun
	return model.FixResult{OK: project.Root != ""}, nil
}
`
	if err := os.WriteFile(filepath.Join(rulesDir, "rules.go"), []byte(rulesSrc), 0644); err != nil {
		t.Fatalf("write rules: %v", err)
	}

	testDir = filepath.Join(moduleRoot, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}

	runGo := `import (
	"testing"

	"` + modPath + `/internal/rules"
	"` + modPath + `/model"
	"github.com/xhd2015/doctest/session"
)

type Request struct {
	// Root becomes model.Project.Root.
	Root string
}

type Response struct {
	OK bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	res, err := rules.FixIgnore(model.Project{Root: req.Root}, false)
	if err != nil {
		return nil, err
	}
	return &Response{OK: res.OK}, nil
}`
	testtree.WriteFile(t, testDir, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))

	fence := string([]byte{'`', '`', '`'})
	writeSubjectLeaf(t, testDir, "leaf-a", fence,
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Root = "/tmp/x"
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
	if resp == nil || !resp.OK {
		t.Fatalf("resp=%v, want OK", resp)
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
