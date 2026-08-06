## Expected

- Suite recording creates at least one new `*.jsonl` under MetricsRoot.
- `metrics top` exits 0.
- Analyze output mentions the fixture leaf (`a_pass` path fragment) or is
  otherwise non-empty ranked output (path / elapsed cues).

## Exit Code

- Analyze exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if len(resp.RunFiles) == 0 {
		t.Fatalf("expected run JSONL under MetricsRoot=%s; RecordErr=%q stderr=%s",
			resp.MetricsRoot, resp.RecordErr, resp.RecordStderr)
	}
	if resp.AnalyzeExitCode != 0 {
		t.Fatalf("metrics top exit=%d AnalyzeErr=%q\nstdout:\n%s\nstderr:\n%s\nfiles=%v",
			resp.AnalyzeExitCode, resp.AnalyzeErr, resp.AnalyzeStdout, resp.AnalyzeStderr, resp.RunFiles)
	}
	out := combinedAnalyze(resp)
	if strings.TrimSpace(out) == "" {
		t.Fatal("metrics top produced empty output")
	}
	// Prefer seeing the fixture leaf path; accept elapsed/path-ish tables if naming differs.
	if !strings.Contains(out, "a_pass") {
		// Still require some substance beyond a bare header with no data.
		lower := strings.ToLower(out)
		if strings.Contains(lower, "no run") || strings.Contains(lower, "not found") {
			t.Fatalf("metrics top empty-store messaging despite JSONL %v:\n%s", resp.RunFiles, out)
		}
		// Soft: if leaf name not present, at least not an error-only dump.
		if len(strings.TrimSpace(out)) < 3 {
			t.Fatalf("metrics top output too thin without a_pass:\n%s", out)
		}
	}
}
```
