# Scenario

**Feature**: `doctest test --go-cmd` selects `go` vs `xgo` by policy and mock detection

```
# parse
doctest test [--go-cmd=auto|xgo|go] <dirs>
  -> Options.GoCmd (default auto when omitted)

# detect (auto only matters for choice)
tree entry packages + project imports
  -> DetectXgoMockUsage
  -> needsXgo if github.com/xhd2015/xgo/runtime/mock reachable transitively

# resolve
ResolveGoTestCmd(mode, needsXgo) -> "go" | "xgo" | error
  auto: needsXgo ? xgo : go
  force: always named binary
  invalid: error

# ensure (when invoking xgo)
EnsureGoTestCmdAvailable("xgo", PATH) -> ok | "xgo not found in PATH"
```

## Preconditions

- Nested self-contained root: does **not** inherit parent `tests/DOCTEST.md` Run.
- Classic TDD: product APIs (`Options.GoCmd`, `DetectXgoMockUsage`,
  `ResolveGoTestCmd`, `EnsureGoTestCmdAvailable`, parse `--go-cmd`) may be
  missing — leaves are expected **RED** until implementer lands them.
- Prefer in-process resolve over real `xgo` / full suite e2e. No trap-all.
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv`.
- Fixture modules use `example.com/app` and stay local (no network).

## Steps

1. Root Setup sets defaults (no parse-only, no check available).
2. Branch Setup sets mode / fixtures.
3. Leaf Setup narrows detection entries, PATH, or parse args.
4. `Run` parse → detect → resolve → optional ensure; Assert checks ResolvedCmd
   or error text.

## Context

- Helpers: `prepareModRoot`, `seedNoMockModule`, `seedTransitiveMockModule`,
  `writeFile`, `fakePATHWithoutXgo`.
- Transitive mock fixture: `runpkg` imports `helper`; `helper` imports
  `github.com/xhd2015/xgo/runtime/mock` (string present in source). Entry import
  path is `example.com/app/runpkg` so detection must follow project imports.
- Fake PATH for missing leaf: temp dir with no `xgo` binary (and preferably no
  accidental xgo from system PATH — searchPATH fully replaces lookup).

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// mockImportPath is the stable import path detection must recognize.
	mockImportPath   = "github.com/xhd2015/xgo/runtime/mock"
	fixtureModPath   = "example.com/app"
	fixtureRunPkg    = "example.com/app/runpkg"
	fixtureHelperPkg = "example.com/app/helper"
)

func errText(resp *Response) string {
	if resp == nil {
		return ""
	}
	return strings.TrimSpace(resp.ErrMsg + "\n" + resp.ParseErr)
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Defaults: resolution path (not parse-only). Leaves override.
	req.ParseOnly = false
	req.CheckAvailable = false
	return nil
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func prepareModRoot(t *testing.T, req *Request) {
	t.Helper()
	req.ModRoot = t.TempDir()
	writeFile(t, filepath.Join(req.ModRoot, "go.mod"),
		"module "+fixtureModPath+"\n\ngo 1.19\n")
}

// seedNoMockModule: entry runpkg has no path to xgo/runtime/mock.
func seedNoMockModule(t *testing.T, req *Request) {
	t.Helper()
	prepareModRoot(t, req)
	writeFile(t, filepath.Join(req.ModRoot, "runpkg", "run.go"),
		"package runpkg\n\n// Tree Run entry: plain helper, no mock.\nfunc Run() string { return \"ok\" }\n")
	req.EntryImportPaths = []string{fixtureRunPkg}
	req.Detect = true
}

// seedTransitiveMockModule: runpkg -> helper -> xgo/runtime/mock (transitive).
// Entry is runpkg only; mock string must NOT be required in runpkg itself.
func seedTransitiveMockModule(t *testing.T, req *Request) {
	t.Helper()
	prepareModRoot(t, req)
	// Blank import is enough: detection keys on import path string, not symbols.
	writeFile(t, filepath.Join(req.ModRoot, "helper", "helper.go"),
		"package helper\n\nimport _ \""+mockImportPath+"\"\n\nfunc Helper() {}\n")
	writeFile(t, filepath.Join(req.ModRoot, "runpkg", "run.go"),
		"package runpkg\n\nimport \""+fixtureHelperPkg+"\"\n\n// Entry calls project helper that pulls xgo mock transitively.\nfunc Run() { helper.Helper() }\n")
	req.EntryImportPaths = []string{fixtureRunPkg}
	req.Detect = true
}

// fakePATHWithoutXgo returns a directory that contains no xgo executable.
// EnsureGoTestCmdAvailable must search only this PATH when SearchPATH is set.
func fakePATHWithoutXgo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// Optional decoy "go" so PATH is non-empty but still lacks xgo.
	goPath := filepath.Join(dir, "go")
	if err := os.WriteFile(goPath, []byte("#!/bin/sh\necho fake-go\n"), 0o755); err != nil {
		t.Fatalf("write fake go: %v", err)
	}
	return dir
}
```
