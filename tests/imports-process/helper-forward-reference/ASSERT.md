---
label: heavy
---

## Expected

- Exit code is 0: doctest generates valid Go from the two helpers and the inner
  leaf test compiles and passes.
- The caller helper is defined before the callee in source order. Top-level funcs
  permit that forward reference, but func literals do not, so the codegen must
  emit the callee's closure before the caller's.
- Stdout/stderr must not contain `undefined:` or `build failed` (the bug symptom
  when helpers are emitted in source order without topological sorting).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	combined := resp.Stdout + "\n" + resp.Stderr
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0 (helper closures must be topologically ordered), got %d\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "undefined:") {
		t.Fatalf("expected no undefined-identifier compile error (forward reference between closures), got:\n%s", combined)
	}
	if strings.Contains(combined, "build failed") {
		t.Fatalf("expected no build failure, got:\n%s", combined)
	}
}
```
