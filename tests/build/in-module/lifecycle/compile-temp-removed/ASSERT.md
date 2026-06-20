## Expected

- `doctest test` exits 0.
- No `.doctest_run_*` directories under `moduleRoot` after run.
- `--gen-dir` dump at `_gen` still contains generated test file.
- Output references `.doctest_run_` temp compile path during run.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected compile temp lifecycle test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoDoctestRunDirs(t, moduleRoot)
	assertDumpHasInternalImport(t, genDir)
	assertStderrUsesTempCompile(t, resp)
}
```