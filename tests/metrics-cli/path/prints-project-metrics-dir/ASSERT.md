## Expected

- Exit code 0.
- Stdout contains the absolute project metrics directory
  `$MetricsRoot/doctest/metrics/<project_id>` (trimmed line acceptable).

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	want, err := filepath.Abs(projectMetricsDir(req))
	if err != nil {
		want = projectMetricsDir(req)
	}
	out := strings.TrimSpace(resp.Stdout)
	if !strings.Contains(out, want) && !strings.Contains(out, projectMetricsDir(req)) {
		t.Fatalf("path stdout missing metrics dir %q:\n%s", want, resp.Stdout)
	}
	mustContain(t, out, req.ProjectID, "path stdout project_id")
	mustContain(t, out, "doctest", "path stdout layout")
	mustContain(t, out, "metrics", "path stdout layout")
}
```
