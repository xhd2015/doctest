## Expected

- Exit non-zero: vacuous Setup (`return nil` only) is rejected.
- stderr tells authors to **remove** the Go **code block** (or omit Setup).
- stderr must **not** say "implement the behavior".
- Message includes a stable marker: `vacuous`, `anti-pattern`, `return nil`, or `blank`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit for vacuous return-nil Setup\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	stderr := resp.Stderr
	if strings.Contains(stderr, "implement the behavior") {
		t.Fatalf("message still asks to implement behavior; should tell author to remove the Go code block:\n%s", stderr)
	}
	if !strings.Contains(stderr, "remove") {
		t.Fatalf("stderr missing remove guidance:\n%s", stderr)
	}
	if !strings.Contains(stderr, "code block") {
		t.Fatalf("stderr missing code block guidance:\n%s", stderr)
	}
	hasMarker := strings.Contains(stderr, "vacuous") ||
		strings.Contains(stderr, "anti-pattern") ||
		strings.Contains(stderr, "return nil") ||
		strings.Contains(stderr, "blank")
	if !hasMarker {
		t.Fatalf("stderr missing vacuous/anti-pattern/return nil/blank marker:\n%s", stderr)
	}
}
```
