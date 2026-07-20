## Expected

- `Run` succeeds (no error).
- Gen root contains `doctest.gen-manifest`.
- Manifest content references `go.mod` (path relative to gen root).
- Manifest includes a version marker (format field so bumps force re-hash miss).
- Generated `go.mod` exists and contains `module testcase`.

## Side Effects

- No requirement on `go.sum` for a parent without sum.
- Legacy `doctest.gomod-fp` must not be the bookkeeping mechanism (see sibling leaf).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	if !resp.ManifestExists {
		t.Fatalf("expected %s after first WriteGoMod; genDir=%s", genManifestName, req.GenDir)
	}
	if resp.ManifestContent == "" {
		t.Fatal("manifest exists but is empty")
	}
	if findManifestLine(resp.ManifestContent, "go.mod") == "" {
		t.Fatalf("manifest must list go.mod, got:\n%s", resp.ManifestContent)
	}
	lower := strings.ToLower(resp.ManifestContent)
	if !strings.Contains(lower, "version") {
		t.Fatalf("manifest must include a version field, got:\n%s", resp.ManifestContent)
	}
	if !strings.Contains(resp.GoModContent, "module testcase") {
		t.Fatalf("expected generated module testcase, got:\n%s", resp.GoModContent)
	}
}
```
