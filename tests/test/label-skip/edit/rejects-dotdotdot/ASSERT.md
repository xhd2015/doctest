## Expected

- Non-zero exit; stderr mentions `...` not accepted.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, "...") {
		t.Fatalf("expected ... mentioned in output\ngot:\n%s", combined)
	}
}
```