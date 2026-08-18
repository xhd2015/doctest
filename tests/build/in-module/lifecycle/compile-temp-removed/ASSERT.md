---
label: e2e
---

## Expected

- `doctest test` exits 0.
- No `.doctest_run_*` directories under `moduleRoot` after run.
- `--gen-dir` dump at `_gen` contains unified gen (expose for parent internal).
- Output uses layout A (suite / mapping-gen), not classic `.doctest_run_`.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected compile temp lifecycle test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoDoctestRunDirs(t, req.ModuleRoot)
	assertDumpHasInternalImport(t, req.GenDir)
	assertStderrUsesTempCompile(t, resp)
}
```