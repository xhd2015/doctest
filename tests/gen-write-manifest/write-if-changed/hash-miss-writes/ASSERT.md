## Expected

- Second write succeeds.
- Target on-disk content reflects the new source (formatted form of sampleGoB).
- Manifest lists the relative path.
- If a pre-change manifest entry existed, it differs after the hash miss.
- No `doctest.gomod-fp`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("hash-miss write failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.TargetContent == "" {
		t.Fatalf("expected target file at %s", req.RelPath)
	}
	// Formatted output may rearrange whitespace/imports; require distinctive token.
	if !strings.Contains(resp.TargetContent, "return 99") {
		t.Fatalf("expected updated body (return 99), got:\n%s", resp.TargetContent)
	}
	if strings.Contains(resp.TargetContent, "return 42") {
		t.Fatalf("stale content still present:\n%s", resp.TargetContent)
	}
	if resp.ManifestEntryAfter == "" {
		t.Fatalf("manifest must list %s after hash-miss write:\n%s", req.RelPath, resp.ManifestContent)
	}
	if req.SnapManifestEntryBefore != "" && resp.ManifestEntryAfter == req.SnapManifestEntryBefore {
		t.Fatalf("expected manifest entry to change on content update:\n%s", resp.ManifestEntryAfter)
	}
}
```
