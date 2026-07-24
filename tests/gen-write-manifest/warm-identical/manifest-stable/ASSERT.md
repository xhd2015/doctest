## Expected

- Second WriteGoMod succeeds.
- `doctest.gen-manifest` exists before and after.
- Manifest mtime equals pre-second-call forced mtime (no rewrite when map unchanged).
- No `doctest.gomod-fp`.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second WriteGoMod failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.ManifestMtimeBefore.IsZero() {
		t.Fatalf("expected manifest present after first write so mtime could be snapshotted (%s)", genManifestName)
	}
	if !resp.ManifestMtimeBefore.Equal(resp.ManifestMtimeAfter) {
		t.Fatalf("manifest rewritten though map unchanged: before=%v after=%v\ncontent:\n%s",
			resp.ManifestMtimeBefore, resp.ManifestMtimeAfter, resp.ManifestContent)
	}
}
```
