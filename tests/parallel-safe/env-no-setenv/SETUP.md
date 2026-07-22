# Scenario

**Feature**: lock zero process env writes for session id and cold GOCACHE while CLI session / cold-cache still work via cmd.Env

```
# static (product contract)
libdoc/**/*.go (non-test)
  -> forbid os/syscall Setenv|Unsetenv of DOCTEST_SESSION_ID / GOCACHE
  -> forbid os.Clearenv in those product paths when paired with the same keys
    (scan focuses on Setenv/Unsetenv for the two keys)

# functional
testbin.Ensure(moduleRoot) -> doctest binary
doctest test <tiny-fixture>                 -> exit 0; nested d.SESSION_ID set
doctest test --cold-cache <tiny-fixture>    -> exit 0; cold GOCACHE via opts→cmd.Env
```

## Preconditions

- Nested root: module root is three levels above `d.DOCTEST_ROOT`
  (`tests/parallel-safe/env-no-setenv` → workspace root).
- Shared doctest binary via `testbin.Ensure` for functional leaves.
- Static leaf does not require the binary.
- Classic TDD: **static leaf RED** until implementer removes `os.Setenv` for
  session/GOCACHE in `libdoc/runner` (and any other product libdoc sites).
- Functional leaves are **GREEN** today and must **stay GREEN** after the fix
  (implementer plumbs `opts` → `cmd.Env` key-replace).
- Out of scope (later phases): `os.Chdir`, package `genDir` vars, Stdout swap.

## Steps

1. Resolve `req.ModuleRoot` from `d.DOCTEST_ROOT`.
2. Build/reuse `req.Bin` for CLI ops.
3. Provide helpers: tiny fixture writers, product source scanner, cold-cache env sandbox.

## Context

- Deep session cookbook: `tests/session-inject/`.
- Deep cold-cache cookbook: `tests/test/cold-cache/`.
- This tree is the P1 anti-pattern + smoke lock only.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

var bt = string([]byte{96, 96, 96}) // triple backtick for nested markdown fixtures

// Matches process env writers that target SESSION_ID or GOCACHE.
// Covers: os.Setenv("DOCTEST_SESSION_ID",…), os.Setenv(core.DoctestSessionIDEnv,…),
// os.Setenv("GOCACHE",…), syscall.Setenv(…), and Unsetenv variants.
var processEnvWriteSessionOrGoCache = regexp.MustCompile(
	`(?m)(?:os|syscall)\.(?:Setenv|Unsetenv)\s*\(\s*(?:[\w.]*DoctestSessionIDEnv|"(?:DOCTEST_SESSION_ID|GOCACHE)")`,
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ModuleRoot = moduleRootFromD(d)
	if req.Timeout == 0 {
		req.Timeout = 30 * time.Second
	}
	if req.Op == "" {
		req.Op = "cli"
	}
	// Functional leaves need the binary; static scan does not.
	if req.Op != "static_scan" {
		req.Bin = testbin.Ensure(t, req.ModuleRoot)
	}
	return nil
}

func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

// scenarioFence wraps a one-line DSN snippet (uses bt so fences are not written literally here).
func scenarioFence(line string) string {
	return bt + "\n" + line + "\n" + bt
}

func tinyRootSetupMD() string {
	return "# Scenario\n\n**Feature**: tiny fixture root\n\n" + scenarioFence("root setup") +
		"\n\n## Steps\n1. no-op\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
}

func tinyLeafSetupMD() string {
	return "# Scenario\n\n**Feature**: tiny leaf\n\n" + scenarioFence("leaf setup") +
		"\n\n## Steps\n1. no-op\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n")
}

// tinyLeafAssertSessionMD requires non-empty d.DOCTEST_SESSION_ID (nested suite inject).
func tinyLeafAssertSessionMD() string {
	code := "import (\n\t\"testing\"\n\n\t\"github.com/xhd2015/doctest/session\"\n)\n\n" +
		"func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {\n" +
		"\t_ = req\n\t_ = resp\n\t_ = err\n" +
		"\tif d == nil || d.DOCTEST_SESSION_ID == \"\" {\n" +
		"\t\tt.Fatal(\"expected non-empty d.DOCTEST_SESSION_ID from suite child env inject\")\n" +
		"\t}\n}\n"
	return "# Scenario\n\n**Feature**: session id present on d\n\n" + scenarioFence("leaf assert session") +
		"\n\n## Expected\n- d.DOCTEST_SESSION_ID non-empty\n\n" +
		doctestGoBlock(code)
}

func tinyLeafAssertOKMD() string {
	return "# Scenario\n\n**Feature**: leaf passes\n\n" + scenarioFence("leaf pass") +
		"\n\n## Expected\n- pass\n\n" +
		doctestGoBlock("import \"testing\"\n\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n")
}

func tinyRunGo() string {
	return `import "testing"

type Request struct{}
type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	_ = req
	return &Response{}, nil
}
`
}

// createTinyFixture writes a one-leaf doctest tree under dir.
// withSessionAssert: leaf Assert checks d.DOCTEST_SESSION_ID non-empty.
func createTinyFixture(t *testing.T, dir string, withSessionAssert bool) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(tinyRunGo())), 0644); err != nil {
		t.Fatalf("write DOCTEST.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(tinyRootSetupMD()), 0644); err != nil {
		t.Fatalf("write root SETUP: %v", err)
	}
	leaf := filepath.Join(dir, "simple")
	if err := os.MkdirAll(leaf, 0755); err != nil {
		t.Fatalf("mkdir leaf: %v", err)
	}
	if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(tinyLeafSetupMD()), 0644); err != nil {
		t.Fatalf("write leaf SETUP: %v", err)
	}
	assertBody := tinyLeafAssertOKMD()
	if withSessionAssert {
		assertBody = tinyLeafAssertSessionMD()
	}
	if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(assertBody), 0644); err != nil {
		t.Fatalf("write leaf ASSERT: %v", err)
	}
}

// createTempModuleFixture returns a module dir containing mytest/ tiny tree.
func createTempModuleFixture(t *testing.T, withSessionAssert bool) (moduleDir, testDir string) {
	t.Helper()
	moduleDir = t.TempDir()
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	testDir = filepath.Join(moduleDir, "mytest")
	createTinyFixture(t, testDir, withSessionAssert)
	return moduleDir, testDir
}

// scanLibdocNoProcessEnvSessionGoCache walks non-test .go under libdoc and
// returns file:line findings for forbidden process env writes targeting
// DOCTEST_SESSION_ID / GOCACHE.
func scanLibdocNoProcessEnvSessionGoCache(moduleRoot string) []string {
	libdoc := filepath.Join(moduleRoot, "libdoc")
	var findings []string
	_ = filepath.WalkDir(libdoc, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			// Skip nested testdata / tests trees under packages if any.
			name := d.Name()
			if name == "testdata" || name == "tests" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		src := string(data)
		// Fast reject: no Setenv/Unsetenv at all.
		if !strings.Contains(src, "Setenv") && !strings.Contains(src, "Unsetenv") {
			return nil
		}
		locs := processEnvWriteSessionOrGoCache.FindAllStringIndex(src, -1)
		if len(locs) == 0 {
			return nil
		}
		rel, relErr := filepath.Rel(moduleRoot, path)
		if relErr != nil {
			rel = path
		}
		for _, loc := range locs {
			line := 1 + strings.Count(src[:loc[0]], "\n")
			snippet := src[loc[0]:loc[1]]
			if len(snippet) > 80 {
				snippet = snippet[:80] + "…"
			}
			findings = append(findings, fmtFinding(rel, line, snippet))
		}
		return nil
	})
	return findings
}

func fmtFinding(rel string, line int, snippet string) string {
	return fmt.Sprintf("%s:%d: %s", strings.TrimSpace(rel), line, strings.TrimSpace(snippet))
}
```
