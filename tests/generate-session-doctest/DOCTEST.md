# Generate pipeline: inject `d *session.Doctest` (Plan P2)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Assembler** — `libdoc/core` generation paths that lower a `TreeCase` into Go source:
  - **classic** — `AssembleTestSource` (single package, everything inlined)
  - **ref** — `AssembleRefRootSource` + `AssembleRefLeafTestSource` (shared root package + thin leaf test)
  - **unified** — `AssembleUnifiedLeafSource` (leaf `RunTestLeaf` package registering with suite)
- **Session Doctest** — public type `session.Doctest` (`github.com/xhd2015/doctest/session`) with fields
  `DOCTEST_ROOT`, `DOCTEST_CASE`, `DOCTEST_SESSION_ID` (struct fields, not free package vars).
- **Author harness** — SETUP / Run / Assert funcs written by test authors. May omit the inject
  param or declare an optional second param after `t` of type `*session.Doctest` (any name).
- **Signature rules** — `libdoc/rules` (+ parse) accept both old (no `d`) and new (with `d`) shapes.
- **Generated test** — the `*_test.go` / leaf package body that constructs `d`, wires Setup→Run→Assert.

**Behaviors**

- At the start of each generated test entry, assembler constructs:
  `d := &session.Doctest{ DOCTEST_ROOT, DOCTEST_CASE, DOCTEST_SESSION_ID }`.
- Call sites always pass that `d` into Setup / Run / Assert.
- If the author omitted the second param, assembler inserts `_ *session.Doctest` in the
  generated signature; if the author named it, the name is preserved.
- No leaf `os.Chdir` / Getwd-restore boilerplate is emitted.
- No package-level free vars `DOCTEST_ROOT` / `DOCTEST_SESSION_ID` (or free-var assignment).
- Import path `github.com/xhd2015/doctest/session` is present when `session.Doctest` is used.
- `DOCTEST_CASE` is the absolute leaf case dir (`filepath.Join(root, rel)` or root when rel empty).

```
TreeCase + author funcs
  -> Assemble{classic|ref|unified}
  -> generated source:
       d := &session.Doctest{...}
       setup(t, d, req) / run(t, d, req) / assert(t, d, req, resp, err)
       # no os.Chdir; no package free DOCTEST_* vars
```

## Decision Tree

```
generate-session-doctest/
├── classic/                              assemble path: AssembleTestSource
│   ├── injects-d/                        d construct + pass; no Chdir; no free vars
│   ├── optional-d-omitted-underscore/    author omits d → `_ *session.Doctest`
│   ├── optional-d-present-keep-name/     author writes `d *session.Doctest` → keep name
│   └── case-path-leaf-abs/               DOCTEST_CASE is abs(root+leaf rel)
├── ref/                                  assemble path: AssembleRefLeafTestSource (+ root)
│   └── injects-d/                        same inject contract on ref leaf (and root free vars gone)
├── unified/                              assemble path: AssembleUnifiedLeafSource
│   └── injects-d/                        same inject contract on unified leaf
└── signature-rules/                      parse/rules accept optional d
    ├── with-d-accepted/                  Setup/Run/Assert with d parse OK
    └── without-d-still-accepted/         classic without-d shapes still parse OK
```

## Test Index

| Leaf | Description |
|------|-------------|
| `classic/injects-d` | Classic source constructs `session.Doctest`, passes `d`, no Chdir, no free DOCTEST_* vars |
| `classic/optional-d-omitted-underscore` | Author omits second param → generated params include `_ *session.Doctest` |
| `classic/optional-d-present-keep-name` | Author uses `d *session.Doctest` → generated keeps that name |
| `classic/case-path-leaf-abs` | `d.DOCTEST_CASE` assignment uses abs path joining root + leaf rel |
| `ref/injects-d` | Ref leaf (and root package) follow inject contract; no Chdir / free vars |
| `unified/injects-d` | Unified leaf follows inject contract; no Chdir / free vars |
| `signature-rules/with-d-accepted` | Parse accepts Setup/Run/Assert signatures that include `d *session.Doctest` |
| `signature-rules/without-d-still-accepted` | Old without-d signatures still accepted |

## How to Run

```sh
doctest vet ./tests/generate-session-doctest/
doctest test -v ./tests/generate-session-doctest/
# expect RED against current assemble (still Chdir + free vars)
```

```go
import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/rules"
)

// Request selects which assemble / parse path Run exercises.
type Request struct {
	// Op is one of:
	//   "classic" | "ref" | "unified" | "parse-with-d" | "parse-without-d"
	Op string

	// AuthorDMode controls second-param after t in author snippets:
	//   "omit" (default) | "named-d" | "named-ctx"
	AuthorDMode string

	// CasePath is the TreeCase.Path (relative leaf under doc root).
	CasePath string

	// DocTestRoot is the absolute root string passed to assemble APIs.
	// Empty → use t.TempDir() absolute path.
	DocTestRoot string
}

// Response carries assembled source text and/or parse outcomes.
type Response struct {
	Source   string
	RootSrc  string // ref root package source when Op=ref
	ParseErr string // non-empty when parse/rules failed
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Op {
	case "classic":
		root := resolveDocRoot(t, req)
		tc := buildTreeCase(req.CasePath, req.AuthorDMode)
		src, err := core.AssembleTestSource(tc, false, "leaf_tc", root)
		if err != nil {
			return nil, err
		}
		return &Response{Source: src}, nil

	case "ref":
		root := resolveDocRoot(t, req)
		tc := buildTreeCase(req.CasePath, req.AuthorDMode)
		rootDocs, _ := core.SplitRefSetupDocs(tc.SetupFiles)
		rootSrc, err := core.AssembleRefRootSource(rootDocs, core.RefRootPkgName)
		if err != nil {
			return nil, err
		}
		leafSrc, err := core.AssembleRefLeafTestSource(tc, false, "leaf_tc", root, core.RefRootImportPath, core.RefRootPkgName)
		if err != nil {
			return nil, err
		}
		return &Response{Source: leafSrc, RootSrc: rootSrc}, nil

	case "unified":
		root := resolveDocRoot(t, req)
		tc := buildTreeCase(req.CasePath, req.AuthorDMode)
		src, err := core.AssembleUnifiedLeafSource(tc, false, "leaf_tc", root, core.RefRootImportPath, core.RefRootPkgName, "testcase/__registry")
		if err != nil {
			return nil, err
		}
		return &Response{Source: src}, nil

	case "parse-with-d":
		return parseSignatures(true)

	case "parse-without-d":
		return parseSignatures(false)

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func resolveDocRoot(t *testing.T, req *Request) string {
	t.Helper()
	if req.DocTestRoot != "" {
		return req.DocTestRoot
	}
	abs, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatalf("abs temp: %v", err)
	}
	return abs
}

func fenceGo(code string) string {
	// Build markdown go fences via char codes so this outer go block is not split.
	bt := "\x60\x60\x60"
	return bt + "go\n" + code + "\n" + bt + "\n"
}

func parseSignatures(withD bool) (*Response, error) {
	setupCode, runCode, assertCode := signatureBodies(withD)

	setupMD := "# Scenario\n\n" + fenceGo(setupCode)
	if _, err := core.ParseSetupDocument("SETUP.md", setupMD); err != nil {
		return &Response{ParseErr: err.Error()}, nil
	}

	// Run lives in DOCTEST.md body.
	doctestMD := "## Version\n0.0.2\n\n# DSN (Domain Specific Notion)\n\nparse fixture\n\n" + fenceGo(runCode)
	if _, err := core.ParseDOCTESTDocument("DOCTEST.md", doctestMD); err != nil {
		return &Response{ParseErr: err.Error()}, nil
	}

	assertMD := "## Expected\n\n" + fenceGo(assertCode)
	if _, err := core.ParseAssertDocument("ASSERT.md", assertMD); err != nil {
		return &Response{ParseErr: err.Error()}, nil
	}

	// Also exercise rules helpers directly (same acceptance).
	setupParams, runParams, assertParams := signatureParams(withD)
	if v := rules.CheckSetupSignature(setupParams, "error", "SETUP.md"); v != nil {
		return &Response{ParseErr: v.Msg}, nil
	}
	if v := rules.CheckRunSignature(runParams, "(*Response, error)", "DOCTEST.md"); v != nil {
		return &Response{ParseErr: v.Msg}, nil
	}
	if v := rules.CheckAssertSignature(assertParams, "", "ASSERT.md"); v != nil {
		return &Response{ParseErr: v.Msg}, nil
	}
	return &Response{ParseErr: ""}, nil
}

func signatureParams(withD bool) (setup, run, assert string) {
	if withD {
		return "t *testing.T, d *session.Doctest, req *Request",
			"t *testing.T, d *session.Doctest, req *Request",
			"t *testing.T, d *session.Doctest, req *Request, resp *Response, err error"
	}
	return "t *testing.T, req *Request",
		"t *testing.T, req *Request",
		"t *testing.T, req *Request, resp *Response, err error"
}

func signatureBodies(withD bool) (setup, run, assert string) {
	imp := `import (
	"testing"

	"github.com/xhd2015/doctest/session"
)
`
	if !withD {
		imp = "import \"testing\"\n"
	}
	sp, rp, ap := signatureParams(withD)
	setup = imp + "func Setup(" + sp + ") error { _ = req; return nil }"
	run = imp + "type Request struct{}\ntype Response struct{}\nfunc Run(" + rp + ") (*Response, error) { return &Response{}, nil }"
	assert = imp + "func Assert(" + ap + ") { _ = req; _ = resp; _ = err }"
	if withD {
		// silence unused d
		setup = imp + "func Setup(" + sp + ") error { _ = d; _ = req; return nil }"
		run = imp + "type Request struct{}\ntype Response struct{}\nfunc Run(" + rp + ") (*Response, error) { _ = d; return &Response{}, nil }"
		assert = imp + "func Assert(" + ap + ") { _ = d; _ = req; _ = resp; _ = err }"
	}
	return setup, run, assert
}

func buildTreeCase(casePath, authorDMode string) core.TreeCase {
	setupParams, runParams, assertParams := authorParams(authorDMode)
	return core.TreeCase{
		Name: "leaf",
		Path: casePath,
		SetupFiles: []core.SetupDocument{
			{
				Path: "DOCTEST.md",
				GoBlock: &core.GoBlock{
					Types: map[string]bool{"Request": true, "Response": true},
					TypeDecls: []string{
						"type Request struct{}",
						"type Response struct{}",
					},
					Run: &core.FuncSnippet{
						Name:        "Run",
						Params:      runParams,
						Results:     "(*Response, error)",
						ResultTypes: "(*Response, error)",
						Body:        "{ return &Response{}, nil }",
					},
					Setup: &core.FuncSnippet{
						Name:    "Setup",
						Params:  setupParams,
						Results: "error",
						Body:    "{ _ = req; return nil }",
					},
				},
			},
			{
				Path: filepath.ToSlash(filepath.Join(casePath, "SETUP.md")),
				GoBlock: &core.GoBlock{
					Setup: &core.FuncSnippet{
						Name:    "Setup",
						Params:  setupParams,
						Results: "error",
						Body:    "{ _ = req; return nil }",
					},
				},
			},
		},
		AssertFile: core.AssertDocument{
			GoBlock: core.GoBlock{
				Assert: &core.FuncSnippet{
					Name:   "Assert",
					Params: assertParams,
					Body:   "{}",
				},
			},
		},
	}
}

func authorParams(mode string) (setup, run, assert string) {
	switch mode {
	case "named-d":
		return "t *testing.T, d *session.Doctest, req *Request",
			"t *testing.T, d *session.Doctest, req *Request",
			"t *testing.T, d *session.Doctest, req *Request, resp *Response, err error"
	case "named-ctx":
		return "t *testing.T, ctx *session.Doctest, req *Request",
			"t *testing.T, ctx *session.Doctest, req *Request",
			"t *testing.T, ctx *session.Doctest, req *Request, resp *Response, err error"
	default: // omit
		return "t *testing.T, req *Request",
			"t *testing.T, req *Request",
			"t *testing.T, req *Request, resp *Response, err error"
	}
}

// --- source contract helpers used by ASSERT leaves ---

func hasSessionDoctestType(src string) bool {
	return strings.Contains(src, "session.Doctest")
}

func hasSessionImport(src string) bool {
	return strings.Contains(src, "github.com/xhd2015/doctest/session")
}

func hasDConstruct(src string) bool {
	// Accept common shapes: d := &session.Doctest{ ... } or d = &session.Doctest{...}
	return strings.Contains(src, "&session.Doctest") &&
		(strings.Contains(src, "d :=") || strings.Contains(src, "d =") || strings.Contains(src, "d:="))
}

func hasLeafChdirBoilerplate(src string) bool {
	// Old generated boilerplate: Getwd + defer Chdir + Chdir into DOCTEST_ROOT / Join(DOCTEST_ROOT, ...)
	if strings.Contains(src, "os.Chdir(filepath.Join(DOCTEST_ROOT") {
		return true
	}
	if strings.Contains(src, "os.Chdir(DOCTEST_ROOT") {
		return true
	}
	if strings.Contains(src, "defer os.Chdir(__origWd)") {
		return true
	}
	if strings.Contains(src, "__origWd, __wdErr := os.Getwd()") {
		return true
	}
	return false
}

func hasPackageFreeDoctestVars(src string) bool {
	// Old package-level free vars (exact spacing from current assembler).
	if strings.Contains(src, "\tDOCTEST_ROOT       string") {
		return true
	}
	if strings.Contains(src, "\tDOCTEST_SESSION_ID string") {
		return true
	}
	// Free-var assignment of the old package vars (not field assign on d / droot).
	if strings.Contains(src, "\tDOCTEST_ROOT = `") || strings.Contains(src, "\tDOCTEST_ROOT = \"") {
		return true
	}
	if strings.Contains(src, "\tDOCTEST_SESSION_ID = sid") {
		return true
	}
	return false
}

func passesDToSetupRunAssert(src string) bool {
	// Call sites must pass the constructed d as second arg after t.
	// setup*(t, d, req) / RootSetup*(t, d, req) / droot.RootSetup*(t, d, req)
	if !strings.Contains(src, "t, d, req") {
		return false
	}
	// run / Run must receive d
	hasRunCall := strings.Contains(src, "run(t, d, req)") ||
		strings.Contains(src, "Run(t, d, req)") ||
		strings.Contains(src, ".Run(t, d, req)")
	if !hasRunCall {
		return false
	}
	// assert / Assert must receive d
	hasAssertCall := strings.Contains(src, "assert(t, d, req") ||
		strings.Contains(src, "Assert(t, d, req")
	return hasAssertCall
}

func hasUnderscoreDoctestParam(src string) bool {
	return strings.Contains(src, "_ *session.Doctest")
}

func hasNamedDoctestParam(src, name string) bool {
	return strings.Contains(src, name+" *session.Doctest")
}

func containsCasePathAssignment(src, absCase string) bool {
	if absCase == "" {
		return false
	}
	// Field assignment inside composite literal or struct field set.
	if !strings.Contains(src, "DOCTEST_CASE") {
		return false
	}
	return strings.Contains(src, absCase)
}
```
