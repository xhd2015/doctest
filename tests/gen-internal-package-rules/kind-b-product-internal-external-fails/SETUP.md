# Scenario

**Feature**: Kind B — product `internal` import under external gen fails

```
module example.com/runner + replace example.com/app => ./app
doctest test ./tests  (leaf imports example.com/app/internal/greet)
  -> CasesImportInternalPackage(runner) false
  -> unified external gen (testcase)
  -> FAIL use of internal package example.com/app/internal/greet
```

## Steps

1. Copy fixture; write runner `go.mod` with replace.
2. `go mod tidy` in runner.
3. Run `doctest test ./tests/...` with WorkDir = runner.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	fixtureRoot := mustCopyFixture(t, d, "product-internal-external-gen")
	runner := t.TempDir()
	// Place app and tests under runner.
	if err := copyDir(filepath.Join(fixtureRoot, "app"), filepath.Join(runner, "app")); err != nil {
		t.Fatalf("copy app: %v", err)
	}
	if err := copyDir(filepath.Join(fixtureRoot, "tests"), filepath.Join(runner, "tests")); err != nil {
		t.Fatalf("copy tests: %v", err)
	}
	gomod := `module example.com/runner

go 1.21

require example.com/app v0.0.0

replace example.com/app => ./app
`
	if err := os.WriteFile(filepath.Join(runner, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatalf("write runner go.mod: %v", err)
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = runner
	tidy.Env = append(os.Environ(), "GOFLAGS=")
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}

	req.Kind = "B"
	req.WorkDir = runner
	req.Args = []string{"test", "-count=1", "-v", "./tests/..."}
	return nil
}
```
