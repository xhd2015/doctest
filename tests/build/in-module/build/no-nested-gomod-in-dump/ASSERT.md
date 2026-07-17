---
label: heavy
---

## Expected

- `doctest build` exits 0.
- Dump at `--gen-dir` contains `tests/leaf/leaf_test.go` importing `internal/greet`.
- Dump has no nested `go.mod`.
- Compile used `.doctest_run_*` temp directory.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected build dump without nested go.mod, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertDumpHasInternalImport(t, genDir)
	assertDumpNoNestedGoMod(t, genDir)
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, moduleRoot)
}
```