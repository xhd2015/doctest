# Scenario

**Feature**: in-repo harnesses consume paths only via `d *session.Doctest`

```
# outer harness (this tree)
generated test -> Setup(t, d, req) / Run(t, d, req) / Assert(t, d, ...)
  d.DOCTEST_ROOT -> module root (../..)
  testbin.Ensure(moduleRoot) -> req.Bin

# fixture under test (temp tree written per leaf)
author Setup/Run/Assert with d
  -> d.DOCTEST_ROOT / d.DOCTEST_CASE
  -> package helpers take d
doctest test <fixture> -> PASS
```

## Preconditions

- This tree is a nested `DOCTEST.md` root under `tests/harness-d-context/` (firewall
  from parent `tests/` free-var harness).
- Module root is two levels above this tree: `filepath.Join(d.DOCTEST_ROOT, "..", "..")`.
- P2 generate already injects `d` and removed free package `DOCTEST_*` vars; fixtures
  that use `d` are expected GREEN without further product changes.
- Subprocess runs with `GOWORK=off` so nested `testcase` modules resolve cleanly.
- Leaves use `e2e` when full integration (otherwise unlabeled) (nested CLI build + `doctest test`).

## Steps

1. Root Setup builds a shared doctest binary via `testbin.Ensure` using `d.DOCTEST_ROOT`.
2. Each leaf writes a temp fixture tree whose author code uses `d *session.Doctest`.
3. Leaf sets `req.Args` to `test <fixture> -v` and runs through root `Run`.
4. Assert expects subprocess exit 0 (fixture PASSes under the inject contract).

## Context

- Outer harness never references free identifiers `DOCTEST_ROOT` / `DOCTEST_SESSION_ID`.
- Helpers that need paths take `d *session.Doctest`.
- Fixture markdown is synthesized in helpers; backticks are built from byte 96.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/session"
)

var bt = string([]byte{96, 96, 96})

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second
	req.Bin = testbin.Ensure(t, moduleRootFrom(d))
	req.Env = append(req.Env, "GOWORK=off")
	return nil
}

// moduleRootFrom resolves the doctest module root from the nested tree's d.DOCTEST_ROOT.
func moduleRootFrom(d *session.Doctest) string {
	return filepath.Join(d.DOCTEST_ROOT, "..", "..")
}

func fenceGo(code string) string {
	return bt + "go\n" + code + "\n" + bt + "\n"
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeScenarioSetup(t *testing.T, path, feature, dsnSnippet, goCode string) {
	t.Helper()
	body := fmt.Sprintf(
		"# Scenario\n\n**Feature**: %s\n\n%s\n%s\n\n## Steps\n1. exercise harness via d\n\n%s",
		feature,
		bt+"\n"+dsnSnippet+"\n"+bt,
		"",
		fenceGo(goCode),
	)
	writeFile(t, path, body)
}

func writeAssert(t *testing.T, path, goCode string) {
	t.Helper()
	body := "## Expected\n\n- fixture PASSes under d inject contract\n\n" + fenceGo(goCode)
	writeFile(t, path, body)
}

func writeDOCTEST(t *testing.T, path, goCode string) {
	t.Helper()
	body := "# Fixture\n\n## Version\n0.0.2\n\n# DSN (Domain Specific Notion)\n\n" +
		"Fixture harness uses d *session.Doctest fields only.\n\n" +
		fenceGo(goCode)
	writeFile(t, path, body)
}

// setFixtureTest configures req to run doctest test against fixtureDir.
func setFixtureTest(req *Request, fixtureDir string) {
	req.FixtureDir = fixtureDir
	req.WorkDir = fixtureDir
	req.Args = []string{"test", fixtureDir, "-v"}
}
```
