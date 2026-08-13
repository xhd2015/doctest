## Expected

- Exit code **0**: generated suite compiles and the inner leaf runs.
- Combined output must **not** contain `jumps over declaration`.

## Side Effects

- Today RED: `go test` fails with
  `goto mock jumps over declaration of cmd` (same as wrk unwind-pipeline CI).
- Desired: generator emits compilable Go for this helper (or rewrites the
  goto/`:=` so the suite builds).

## Errors

- None on the success path.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	if resp.ExitCode != 0 {
		t.Fatalf("expected inner doctest test to compile and run (exit 0); exit=%d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(out, "jumps over declaration") {
		t.Fatalf("generated suite must not fail with goto-over-declaration\n%s", out)
	}
}
```
