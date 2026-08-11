## Expected

- Subject run succeeds (`RunErr` empty).
- Gen dir contains unified markers: `suite` and `__allleaves`.
- Preferably also `__droot` and `__registry` (layout A full set).

Today (pre-P2): internalCompile classic dump under GenDir lacks suite/__allleaves
— **RED**.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.RunErr != "" {
		t.Fatalf("P2 RED: want subject PASS under unified gen, got RunErr=%s\nstderr:\n%s\ngen=%s",
			resp.RunErr, resp.Stderr, resp.GenDir)
	}
	if resp.GenDir == "" {
		t.Fatal("GenDir empty; cannot inspect unified layout")
	}
	if !resp.HasSuite {
		t.Fatalf("P2 RED: gen missing suite package (unified layout A); go files=%v gen=%s",
			basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasAllLeaves {
		t.Fatalf("P2 RED: gen missing __allleaves (unified layout A); go files=%v gen=%s",
			basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasDroot {
		t.Fatalf("P2 RED: gen missing __droot; go files=%v gen=%s", basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasRegistry {
		t.Fatalf("P2 RED: gen missing __registry; go files=%v gen=%s", basenames(resp.GoFiles), resp.GenDir)
	}
}

func basenames(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = filepath.Base(p)
	}
	return out
}
```
