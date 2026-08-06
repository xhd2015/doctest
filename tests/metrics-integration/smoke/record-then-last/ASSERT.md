## Expected

- Suite recording creates at least one new `*.jsonl` under MetricsRoot.
- Recording does not leave a fatal harness error that prevents any file (prefer
  `RecordErr` empty when possible; file presence is the hard requirement).
- `metrics last` exits 0.
- Analyze combined output is non-empty and does not look like a pure “no runs”
  failure (should not be the only story when a file exists).

## Exit Code

- Analyze exit code 0.

## Side Effects

- One or more run JSONL files under `$MetricsRoot/doctest/metrics/.../runs/`.

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
		t.Fatalf("expected run JSONL under MetricsRoot=%s; RecordErr=%q stderr=%s stdout=%s",
			resp.MetricsRoot, resp.RecordErr, resp.RecordStderr, resp.RecordStdout)
	}
	if resp.AnalyzeExitCode != 0 {
		t.Fatalf("metrics last exit=%d AnalyzeErr=%q\nstdout:\n%s\nstderr:\n%s\nrun files: %v",
			resp.AnalyzeExitCode, resp.AnalyzeErr, resp.AnalyzeStdout, resp.AnalyzeStderr, resp.RunFiles)
	}
	out := combinedAnalyze(resp)
	if strings.TrimSpace(out) == "" {
		t.Fatal("metrics last produced empty output")
	}
	// Soft markers: either counts, path cues, or run-ish tokens — reject pure empty-store messaging alone.
	lower := strings.ToLower(out)
	noRunsOnly := (strings.Contains(lower, "no run") || strings.Contains(lower, "no metrics") || strings.Contains(lower, "not found")) &&
		!strings.Contains(lower, "pass") &&
		!strings.Contains(lower, "total") &&
		!strings.Contains(out, "a_pass")
	if noRunsOnly {
		t.Fatalf("metrics last looks like empty store despite JSONL %v:\n%s", resp.RunFiles, out)
	}
}
```
