## Expected

- Exit code 0 (path is informational; missing dir is not an error).
- Stdout still includes the canonical `…/doctest/metrics/<project_id>` path.

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
		t.Fatalf("expected exit 0 when metrics dir missing; exit=%d stderr=%s", resp.ExitCode, resp.Stderr)
	}
	out := strings.TrimSpace(resp.Stdout)
	mustContain(t, out, req.ProjectID, "path when missing")
	mustContain(t, out, "metrics", "path when missing")
}
```
