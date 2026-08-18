## Expected

**Desired product behavior** (GREEN after fix; **RED** on current code):

- `RunErr` empty: nested subject tree completes under cover + coverpkg.
- Cover profile at `req.CoverPath` exists and is non-empty.
- Combined stdout/stderr do **not** contain the classic expose cover failure:
  `cover:` + `__doctest_internal_expose` + `no such file or directory`.
- Prefer single suite package args when the go test line is present.

Today: expose is overlay-only; `go tool cover` opens the logical
`…/__doctest_internal_expose/…/expose.go` path on disk → open fails → **RED**.

```go
import (
	"os"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if req.CoverPath == "" {
		t.Fatal("CoverPath empty")
	}
	if req.CoverPkg == "" {
		t.Fatal("CoverPkg empty — this leaf requires product-module coverpkg")
	}

	combined := resp.Stdout + "\n" + resp.Stderr + "\n" + resp.RunErr
	exposeCoverOpen := strings.Contains(combined, "__doctest_internal_expose") &&
		strings.Contains(combined, "no such file or directory") &&
		(strings.Contains(combined, "cover:") || strings.Contains(combined, "expose.go"))

	if resp.RunErr != "" {
		hint := ""
		if exposeCoverOpen {
			hint = " (expose: go tool cover opens overlay-only expose.go)"
		}
		t.Fatalf("want cover+coverpkg success for parent-internal expose, got RunErr=%s%s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, hint, resp.Stdout, resp.Stderr)
	}
	if exposeCoverOpen {
		t.Fatalf("cover must not open missing __doctest_internal_expose.go under coverpkg=%q\nstdout:\n%s\nstderr:\n%s",
			req.CoverPkg, resp.Stdout, resp.Stderr)
	}

	info, statErr := os.Stat(req.CoverPath)
	if statErr != nil {
		t.Fatalf("expected coverprofile at %s: %v\nstdout:\n%s\nstderr:\n%s",
			req.CoverPath, statErr, resp.Stdout, resp.Stderr)
	}
	if info.Size() == 0 {
		t.Fatalf("coverprofile %s is empty", req.CoverPath)
	}
	if len(resp.GoTestPackageArgs) > 1 {
		t.Fatalf("coverpkg path should be single-package suite run, got pkgs=%v line=%q",
			resp.GoTestPackageArgs, resp.GoTestDisplayLine)
	}
	if len(resp.GoTestPackageArgs) == 1 && !strings.Contains(resp.GoTestPackageArgs[0], "suite") {
		t.Fatalf("single package should be suite, got %q", resp.GoTestPackageArgs[0])
	}
}
```
