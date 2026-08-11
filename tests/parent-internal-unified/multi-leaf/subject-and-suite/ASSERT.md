## Expected

- `RunErr` empty: both subject leaves (`leaf-a` Hello / `leaf-b` DefaultName) pass.
- Displayed `go test` line exists.
- Package path args have **exactly one** entry containing `suite`.
- Package args are **not** multi-leaf `./leaf-a` + `./leaf-b` (classic internalCompile).

Today (pre-P2): multi-leaf `internalCompile` may still run subject tests but
lists multiple packages — **RED** on suite-only assert.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.RunErr != "" {
		t.Fatalf("P2 RED: want subject multi-leaf parent-internal PASS under unified, got RunErr=%s\nstdout:\n%s\nstderr:\n%s\ngen=%s",
			resp.RunErr, resp.Stdout, resp.Stderr, resp.GenDir)
	}
	if resp.GoTestDisplayLine == "" {
		t.Fatalf("no go test display line in stdout/stderr:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	pkgs := resp.GoTestPackageArgs
	if len(pkgs) != 1 {
		t.Fatalf("P2 RED: want exactly 1 go test package (suite), got %d args=%v line=%q\n(classic internalCompile multi-leaf lists ./leaf-a ./leaf-b)",
			len(pkgs), pkgs, resp.GoTestDisplayLine)
	}
	if !strings.Contains(pkgs[0], "suite") {
		t.Fatalf("single package arg must refer to suite, got %q line=%q", pkgs[0], resp.GoTestDisplayLine)
	}
	joined := strings.Join(pkgs, " ")
	if strings.Contains(joined, "leaf-a") && strings.Contains(joined, "leaf-b") {
		t.Fatalf("go test still multi-package leaf shape: %v line=%q", pkgs, resp.GoTestDisplayLine)
	}
}
```
