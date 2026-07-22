## Expected

- Under `req.ModuleRoot/libdoc`, every non-test `.go` file (excluding package
  `testdata/` and `tests/` dirs) has **zero** process env writes that target
  `DOCTEST_SESSION_ID` / `DoctestSessionIDEnv` or `GOCACHE` via:
  - `os.Setenv` / `syscall.Setenv`
  - `os.Unsetenv` / `syscall.Unsetenv`
- Findings report `path:line: snippet` for implementer fix sites.
- Classic TDD: **RED** while `libdoc/runner/runner.go` still calls
  `os.Setenv` for session id and GOCACHE.

## Side Effects

- None (read-only source scan).

## Exit Code

- N/A for static_scan Run (no-op). Assert fails the test on findings.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = resp
	_ = err
	if req.ModuleRoot == "" {
		t.Fatal("req.ModuleRoot empty")
	}
	findings := scanLibdocNoProcessEnvSessionGoCache(req.ModuleRoot)
	if len(findings) == 0 {
		return
	}
	t.Fatalf("product libdoc must not process-Setenv/Unsetenv DOCTEST_SESSION_ID or GOCACHE (P1); %d finding(s):\n  %s\n\nImplementer: store sid/GoCache on opts and pass via cmd.Env key-replace only.",
		len(findings), strings.Join(findings, "\n  "))
}
```
