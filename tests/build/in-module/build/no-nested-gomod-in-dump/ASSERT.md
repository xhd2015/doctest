---
label: e2e
---

## Expected

- `doctest build` exits 0.
- Dump at `--gen-dir` is unified layout A (`tests/__droot` / leaf) with Kind B
  expose import for parent `internal/greet` (not classic multi-leaf `_test.go`).
- Dump is module `testcase` (go.mod with replace to parent).
- No classic `.doctest_run_*` under the product module root.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected build dump without nested go.mod, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertDumpHasInternalImport(t, req.GenDir)
	assertDumpNoNestedGoMod(t, req.GenDir)
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, req.ModuleRoot)
}
```