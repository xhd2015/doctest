## Expected

- `doctest test` exits 0 via internal-compile (unchanged from pre-assert behavior).
- Verbose output does not reference assert cache replace path.
- No `.doctest_run_*` dirs remain.

## Exit Code

- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected internal-only test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if strings.Contains(combined, "replace "+assertModPath) {
		t.Fatalf("expected no assert replace in internal-only run, got:\n%s", combined)
	}
	assertStderrUsesTempCompile(t, resp)
	assertNoDoctestRunDirs(t, moduleRoot)
}
```