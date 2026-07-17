---
label: heavy
---

## Expected

- `doctest test` exits 0.
- Dump at `--gen-dir` contains `tests/leaf/leaf_test.go` importing `internal/greet`.
- Dump has no nested `go.mod`.
- Compile used `.doctest_run_*` temp (not gen-dir as compile root).
- No `.doctest_run_*` dirs remain under `moduleRoot` after run.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected internal import with gen-dir dump to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertDumpHasInternalImport(t, genDir)
	assertDumpNoNestedGoMod(t, genDir)
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, moduleRoot)
}
```