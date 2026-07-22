---
label: heavy
---

## Expected

- `doctest test` exits 0 via internal-compile.
- Verbose output references `.doctest_run_` temp compile dir.
- Verbose output includes `-modfile=` (always-on assert+session on internal path).
- No `.doctest_run_*` dirs remain after run.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected internal-only test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertStderrUsesTempCompile(t, resp)
	assertStderrUsesModfile(t, resp)
	assertNoDoctestRunDirs(t, req.ModuleRoot)
}
```
