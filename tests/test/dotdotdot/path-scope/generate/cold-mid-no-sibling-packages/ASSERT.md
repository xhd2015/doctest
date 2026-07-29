---
label: heavy
---

## Expected

Cold mid path generate writes mid packages only — no `tree/sibling` under gen.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	out := pathScopeOut(resp)
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, out)
	}
	genDir := filepath.Join(req.WorkDir, ".gen")
	midLeaf := filepath.Join(genDir, "tree", "mid", "two", "leaf.go")
	if _, err := os.Stat(midLeaf); err != nil {
		t.Fatalf("expected mid leaf under gen: %v\n%s", err, out)
	}
	sib := filepath.Join(genDir, "tree", "sibling")
	if st, err := os.Stat(sib); err == nil && st.IsDir() {
		t.Fatalf("cold mid gen must not create sibling packages under %s", sib)
	}
	_ = filepath.Walk(genDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(genDir, p)
		rel = filepath.ToSlash(rel)
		if strings.Contains(rel, "tree/sibling") {
			t.Errorf("out-of-scope sibling gen path: %s", rel)
		}
		return nil
	})
}
```
