---
label: e2e, heavy
---

## Expected

- `doctest test` exits 0 via temp compile under parent module.
- No `.doctest_run_*` directories remain under `moduleRoot` after the run.

## Side Effects

- Compile temp is created and removed; no persistent gen-dir dump.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected internal import test to pass via temp compile, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoDoctestRunDirs(t, req.ModuleRoot)
}
```