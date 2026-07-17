---
label: heavy
---

## Expected

- `doctest test` exits 0 despite gen-dir being outside the parent module.
- Outside dump contains generated test importing `internal/greet`.
- Outside dump has no nested `go.mod`.
- Compile temp under `moduleRoot` is removed after run.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected outside gen-dir dump to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertDumpHasInternalImport(t, outsideGenDir)
	assertDumpNoNestedGoMod(t, outsideGenDir)
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, moduleRoot)
}
```