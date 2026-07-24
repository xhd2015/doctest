## Expected

- Exit code 0.
- stdout contains `--metrics-on`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(out, "--metrics-on") {
		t.Fatalf("test --help missing --metrics-on:\n%s", out)
	}
}
```
