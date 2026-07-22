# Scenario

**Feature**: no `os.Setenv` / `syscall.Setenv` / `Unsetenv` of DOCTEST_SESSION_ID or GOCACHE in product libdoc

```
# inventory to kill (P1)
libdoc/runner/runner.go  Setenv SESSION_ID  -> remove; opts.SessionID + cmd.Env
libdoc/runner/runner.go  Setenv GOCACHE     -> remove; opts.GoCache + cmd.Env
# any other libdoc product sites with the same pattern also forbidden
```

## Preconditions

- Parent set `req.Op=static_scan` and `req.ModuleRoot`.
- Unit-test files (`*_test.go`) are out of scope for this scan (rewritten separately).

## Steps

1. Inherit static_scan op.
2. Assert scans and fails with file:line if any finding remains.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Parent already set static_scan + ModuleRoot.
	if req.ModuleRoot == "" {
		t.Fatal("req.ModuleRoot must be set by root Setup")
	}
	if req.Op != "static_scan" {
		t.Fatalf("expected Op=static_scan, got %q", req.Op)
	}
	return nil
}
```
