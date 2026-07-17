## Expected

- Exit code 0.
- stdout (trimmed) is valid JSON.
- JSON text includes `group/slow-leaf` and an elapsed field (elapsed_ns or similar).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	if !isValidJSON(resp.Stdout) {
		t.Fatalf("top --json stdout not valid JSON:\n%s", resp.Stdout)
	}
	mustContain(t, resp.Stdout, "group/slow-leaf", "top json path")
	low := strings.ToLower(resp.Stdout)
	if !strings.Contains(low, "elapsed") {
		t.Fatalf("top --json should include elapsed timing field:\n%s", resp.Stdout)
	}
}
```
