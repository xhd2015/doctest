## Expected

- Exit code 0.
- stdout contains `--no-metrics`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(out, "--no-metrics") {
		t.Fatalf("test --help missing --no-metrics:\n%s", out)
	}
}
```
