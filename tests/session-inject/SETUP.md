# Scenario

**Feature**: doctest materializes embedded session as a cached local module and replace

```
# session import detected in SETUP/ASSERT Go blocks
doctest test/build <tree>
  -> CasesImportSessionPackage
  -> MaterializeSessionModule
  -> $UserCacheDir/doctest/session-mod/<md5>/{go.mod,*.go}

# nested module replace
session import -> replace github.com/xhd2015/doctest/session => <cache> in testcase go.mod

# consumer Once
leaf imports session -> Once(t, key, fn) -> json.RawMessage
```

## Preconditions

- Module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf builds a fresh doctest binary and runs it in a temp Go module.
- `GOWORK=off` for subprocess invocations.
- Session cache lives at `$UserCacheDir/doctest/session-mod/<content-md5>/`.
- Classic TDD: expect RED until `sessionmod`, materialize, and import detection exist.
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv` in harness helpers that
  intend product semantics; use injected `DOCTEST_SESSION_ID` only for this tree's
  own locks if needed.

## Steps

1. Build the doctest binary from the module root via `testbin.Ensure`.
2. Create a temp module and doctest tree per leaf scenario.
3. Execute the doctest binary with leaf-specific args and inspect outputs/cache.

## Context

- Helpers mirror `tests/embed-assert` but target session import path and
  `session-mod` cache.
- Cache key comes from `sessionmod.RawSourceCacheKeyMD5()` once that package lands.
- Until implementation exists, Setup may fail at import/build — that is the RED state.
- Fixture paths (`ModuleRoot`, `TestDir`) live on `req` (request-local; Parallel-safe).

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/sessionmod"
	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

const (
	modPath        = "example.com/session-app"
	sessionModPath = "github.com/xhd2015/doctest/session"
)

var (
	// Immutable fence marker only — not mutated across leaves.
	bt = string([]byte{96, 96, 96})
)

func lockCacheTests(t *testing.T) {
	t.Helper()
	lockPath := filepath.Join(os.TempDir(), "doctest-session-inject-cache-tests.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open cache test lock: %v", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		t.Fatalf("acquire cache test lock: %v", err)
	}
	t.Cleanup(func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	})
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second
	// Default L2: no binary. Cache leaves set UseCLI+Bin.
	return nil
}

func expectedSessionCacheDir(t *testing.T) string {
	t.Helper()
	base, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("UserCacheDir: %v", err)
	}
	md5hex := sessionmod.RawSourceCacheKeyMD5()
	return filepath.Join(base, "doctest", "session-mod", md5hex)
}

func assertSessionCacheLayout(t *testing.T, cacheDir string) {
	t.Helper()
	goMod := filepath.Join(cacheDir, "go.mod")
	data, err := os.ReadFile(goMod)
	if err != nil {
		t.Fatalf("read session-mod go.mod: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "module "+sessionModPath) {
		t.Fatalf("go.mod missing module %s:\n%s", sessionModPath, text)
	}
	// At least one .go source must exist (once.go / session.go / etc.).
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatalf("readdir cache: %v", err)
	}
	hasGo := false
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") && e.Name() != "go.mod" {
			hasGo = true
			break
		}
	}
	if !hasGo {
		t.Fatalf("session-mod cache %s has no .go sources", cacheDir)
	}
}

func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func createDoctestRoot(dir string, extraImports string, runCode string) error {
	goBody := "import (\n" +
		"\t\"testing\"\n" +
		extraImports +
		")\n\n" +
		"type Request struct{}\n" +
		"type Response struct{ Message string }\n\n" +
		runCode
	return os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(goBody)), 0644)
}

func createDoctestLeaf(dir string, setupGo string, assertGo string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if setupGo == "" {
		setupGo = "import \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }"
	}
	if assertGo == "" {
		assertGo = "import \"testing\"\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {\n" +
			"\tif err != nil { t.Fatal(err) }\n" +
			"\tif resp.Message != \"ok\" { t.Fatalf(\"expected ok, got %q\", resp.Message) }\n" +
			"}"
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(doctestGoBlock(setupGo)), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(doctestGoBlock(assertGo)), 0644)
}

func createPublicModuleProject(t *testing.T, req *Request, leafSetupGo string, leafAssertGo string, withSessionImport bool) {
	t.Helper()

	req.ModuleRoot = t.TempDir()
	if err := os.WriteFile(
		filepath.Join(req.ModuleRoot, "go.mod"),
		[]byte("module "+modPath+"\n\ngo 1.21\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	req.TestDir = filepath.Join(req.ModuleRoot, "tests")
	if err := os.MkdirAll(req.TestDir, 0755); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}

	extraImports := ""
	runCode := "func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {\n" +
		"\treturn &Response{Message: \"ok\"}, nil\n" +
		"}"
	if withSessionImport {
		extraImports = "\t\"encoding/json\"\n\t\"" + sessionModPath + "\"\n"
		// Run imports session so assemble detects the import even if leaf also imports it.
		runCode = "func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {\n" +
			"\t_ = json.RawMessage(nil)\n" +
			"\t_ = session.DoctestSessionIDEnv\n" +
			"\treturn &Response{Message: \"ok\"}, nil\n" +
			"}"
	}
	if err := createDoctestRoot(req.TestDir, extraImports, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(req.TestDir, "leaf"), leafSetupGo, leafAssertGo); err != nil {
		t.Fatalf("create doctest leaf: %v", err)
	}
}

func setupModuleEnv(t *testing.T, req *Request) {
	t.Helper()
	// Parallel-safe: absolute paths in Args; no WorkDir/Env.
	req.WorkDir = ""
}

func defaultSessionAssertGo() string {
	return "import (\n" +
		"\t\"encoding/json\"\n" +
		"\t\"testing\"\n" +
		"\t\"" + sessionModPath + "\"\n" +
		")\n" +
		"func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {\n" +
		"\tif err != nil { t.Fatal(err) }\n" +
		"\traw, e := session.Once(t, \"probe\", func(t testing.TB, cacheDir string) (json.RawMessage, error) {\n" +
		"\t\treturn json.RawMessage(`{\"path\":\"x\"}`), nil\n" +
		"\t})\n" +
		"\tif e != nil { t.Fatal(e) }\n" +
		"\tif !json.Valid(raw) { t.Fatalf(\"invalid json %s\", raw) }\n" +
		"\tif resp.Message != \"ok\" { t.Fatalf(\"expected ok, got %q\", resp.Message) }\n" +
		"}"
}

func findNestedGoMod(t *testing.T, moduleRoot string) string {
	t.Helper()
	var found string
	_ = filepath.WalkDir(moduleRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if d.Name() == "go.mod" && path != filepath.Join(moduleRoot, "go.mod") {
			found = path
			return fmt.Errorf("stop")
		}
		return nil
	})
	return found
}
```
