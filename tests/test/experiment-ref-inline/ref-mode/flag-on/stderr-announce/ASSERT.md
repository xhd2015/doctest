## Expected

- Suite run succeeds.
- Combined stderr (and stdout fallback) contains `experiment` and
  `ref-instead-of-inline` (case-insensitive), documenting experimental mode.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("ref-mode RunTest failed: %s\nstderr:\n%s", resp.RunErr, resp.Stderr)
	}
	out := strings.ToLower(resp.Stderr + "\n" + resp.Stdout)
	if !strings.Contains(out, "experiment") {
		t.Fatalf("stderr/stdout missing %q (ref-mode announce):\nstderr:\n%s\nstdout:\n%s",
			"experiment", resp.Stderr, resp.Stdout)
	}
	if !strings.Contains(out, "ref-instead-of-inline") {
		t.Fatalf("stderr/stdout missing %q (ref-mode announce):\nstderr:\n%s\nstdout:\n%s",
			"ref-instead-of-inline", resp.Stderr, resp.Stdout)
	}
}
```
