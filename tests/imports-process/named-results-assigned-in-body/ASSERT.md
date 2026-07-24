---
label: heavy
---

## Expected

- Exit code is 0: doctest generates valid Go from the named-result helper and the
  inner leaf test compiles and passes.
- The helper body assigns to named result variables (`mainRepo`, `wtDir`,
  `branch`), so the closure signature must preserve those names. Stripping them
  to type-only `(string, string, string)` leaves the assignments referencing
  undeclared identifiers and breaks compilation.
- Stdout/stderr must not contain `undefined:` or `build failed` (the bug symptom
  when `writeFuncClosure` strips named results to type-only).

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
		t.Fatalf("expected exit code 0 (named results must be preserved so the body compiles), got %d\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "undefined:") {
		t.Fatalf("expected no undefined-identifier compile error (named results stripped), got:\n%s", combined)
	}
	if strings.Contains(combined, "build failed") {
		t.Fatalf("expected no build failure, got:\n%s", combined)
	}
}
```
