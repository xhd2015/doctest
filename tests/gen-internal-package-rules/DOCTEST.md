# Gen internal package rules — RED contracts (kind A + kind B)

## Version

0.0.1

## Layer

**L3 e2e, heavy** — builds product `doctest` binary and runs it on fixture
subject trees. Outer leaves expect **FAIL** until gen packaging is fixed
(kind A shim) or until the subject is run under product-module internal-compile
(kind B success path lives in `tests/build/in-module/`).

## DSN

Two illegal-import modes (MECE):

| Kind | Illegal edge | Fixture |
|------|--------------|---------|
| **A** | `__allleaves` → gen path `…/http/internal/…` | `testdata/scenario-path-internal` |
| **B** | gen package → product `example.com/app/internal/…` under external gen | `testdata/product-internal-external-gen` |

Kind A uses **public/stdlib-only** leaf `Run` (no product internal import).
Kind B uses a **runner** module (`example.com/runner`) + `replace` to `app`
so `CasesImportInternalPackage` does not trigger in-module compile.

## Decision tree

```text
gen-internal-package-rules/
├── testdata/scenario-path-internal/          # kind A subject
├── testdata/product-internal-external-gen/   # kind B subject
├── kind-a-scenario-path-fails/               # expect FAIL
└── kind-b-product-internal-external-fails/   # expect FAIL
```

## How to run

```sh
go run ./cmd/doctest test --label 'e2e && heavy' -v \
  ./tests/gen-internal-package-rules/...
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
	// Scenario for Assert diagnostics only.
	Kind string // "A" | "B"
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Combined string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req.Bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	if req.Timeout <= 0 {
		req.Timeout = 120 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	resp := &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Combined: stdout.String() + stderr.String(),
		Err:      err,
	}
	if err == nil {
		resp.ExitCode = 0
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	if ctx.Err() != nil {
		return resp, ctx.Err()
	}
	resp.ExitCode = 1
	return resp, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(dst, 0755)
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

func mustCopyFixture(t *testing.T, d *session.Doctest, rel string) string {
	t.Helper()
	src := filepath.Join(d.DOCTEST_ROOT, "testdata", rel)
	dst := filepath.Join(t.TempDir(), rel)
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("copy fixture %s: %v", rel, err)
	}
	return dst
}

func containsInternalDenied(s string) bool {
	return strings.Contains(s, "use of internal package") ||
		(strings.Contains(s, "internal package") && strings.Contains(s, "not allowed"))
}

// discard unused import guard for io in case of future helpers
var _ = io.Discard
```
