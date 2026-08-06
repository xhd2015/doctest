## Expected

- Parse succeeds.
- `Response.Overlay` is absolute (`filepath.IsAbs`).
- `Response.Overlay` equals `filepath.Abs` of the relative arg (allowing
  macOS `/private/var` vs `/var` normalization used by profile abs helpers).
- Relative path string must not equal `Response.Overlay` (must have been resolved).

## Errors

- Classic TDD **RED** until flag accept + abs-resolve land.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected parse success, got exit=%d err=%s", resp.ExitCode, errText(resp))
	}
	var rel string
	for i, a := range req.ParseArgs {
		if a == "-overlay" && i+1 < len(req.ParseArgs) {
			rel = req.ParseArgs[i+1]
			break
		}
	}
	if rel == "" || filepath.IsAbs(rel) {
		t.Fatalf("SETUP must pass relative -overlay path, got %q", rel)
	}
	if !absLike(resp.Overlay) {
		t.Fatalf("Overlay must be absolute, got %q", resp.Overlay)
	}
	if resp.Overlay == rel {
		t.Fatalf("relative path was not resolved: %q", resp.Overlay)
	}
	want, aerr := filepath.Abs(rel)
	if aerr != nil {
		t.Fatalf("abs: %v", aerr)
	}
	// Accept /private/var ↔ /var equivalence (see absProfilePath).
	got := resp.Overlay
	if got != want && stripPrivateVar(got) != stripPrivateVar(want) {
		t.Fatalf("Overlay=%q want Abs(%q)=%q", resp.Overlay, rel, want)
	}
}

func stripPrivateVar(p string) string {
	const priv = "/private"
	if strings.HasPrefix(p, priv+"/var/") {
		return p[len(priv):]
	}
	return p
}
```
