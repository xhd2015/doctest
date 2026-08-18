---
label: e2e
---

## Expected

- `doctest test` exits 0 via **unified layout A** (mapping-gen + expose).
- Verbose output shows suite / mapping-gen (not classic `-modfile=` / `.doctest_run_*`).
- No `.doctest_run_*` dirs remain under the consumer module root.
- No leftover `.doctest.mod` / `.doctest.sum` under the consumer module root.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected internal+assert modfile test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertStderrUsesTempCompile(t, resp)
	assertStderrUsesModfile(t, resp)
	assertNoDoctestRunDirs(t, req.ModuleRoot)
	assertNoInternalModfileArtifacts(t, req.ModuleRoot)
}
```