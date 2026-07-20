## Expected

- Second write succeeds.
- Target file mtime equals forced pre-second mtime (no rewrite on hash hit).
- `doctest.gen-manifest` lists the relative path.
- No `doctest.gomod-fp`.

## Side Effects

- Manifest map may be unchanged → manifest file mtime may also stay stable
  (not required here; covered under warm-identical/manifest-stable for go.mod).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second write failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.TargetMtimeBefore.IsZero() {
		t.Fatal("missing target mtime snapshot")
	}
	if !resp.TargetMtimeBefore.Equal(resp.TargetMtimeAfter) {
		t.Fatalf("target rewritten on hash hit: before=%v after=%v path=%s",
			resp.TargetMtimeBefore, resp.TargetMtimeAfter, req.RelPath)
	}
	if resp.ManifestEntryAfter == "" {
		t.Fatalf("manifest must list %s after WriteIfChanged, got:\n%s", req.RelPath, resp.ManifestContent)
	}
}
```
