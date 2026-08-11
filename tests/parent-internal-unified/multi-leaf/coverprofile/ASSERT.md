## Expected

- `RunErr` empty (coverprofile accepted for single suite package).
- Cover profile file at `req.CoverPath` exists and is non-empty.
- Preferably go test package args are a single suite package.

Today (pre-P2): multi-leaf internalCompile rejects `-coverprofile` with
multiple packages — **RED**.

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
	if resp.RunErr != "" {
		// Help diagnose classic multi-package coverprofile rejection while RED.
		hint := ""
		if strings.Contains(resp.RunErr, "coverprofile") && strings.Contains(resp.RunErr, "multiple packages") {
			hint = " (classic internalCompile multi-leaf blocks -coverprofile)"
		}
		t.Fatalf("P2 RED: want coverprofile success for parent-internal multi-leaf, got RunErr=%s%s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, hint, resp.Stdout, resp.Stderr)
	}
	info, statErr := os.Stat(req.CoverPath)
	if statErr != nil {
		t.Fatalf("expected coverprofile at %s: %v\nstdout:\n%s\nstderr:\n%s",
			req.CoverPath, statErr, resp.Stdout, resp.Stderr)
	}
	if info.Size() == 0 {
		t.Fatalf("coverprofile %s is empty", req.CoverPath)
	}
	// Soft packaging check when go test line is present.
	if len(resp.GoTestPackageArgs) > 1 {
		t.Fatalf("coverprofile path should be single-package suite run, got pkgs=%v line=%q",
			resp.GoTestPackageArgs, resp.GoTestDisplayLine)
	}
	if len(resp.GoTestPackageArgs) == 1 && !strings.Contains(resp.GoTestPackageArgs[0], "suite") {
		t.Fatalf("single package should be suite, got %q", resp.GoTestPackageArgs[0])
	}
}
```
