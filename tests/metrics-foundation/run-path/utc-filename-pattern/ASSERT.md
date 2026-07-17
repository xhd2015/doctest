## Expected

- Path ends with `doctest/metrics/github.com_xhd2015_doctest/runs/2026-07-17-12-34-56-07-abc12def.jsonl`.
- Basename matches `^\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}-[0-9a-fA-F]+\.jsonl$` style (here suffix is hex-like).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantSuffix := filepath.Join(
		"doctest", "metrics", "github.com_xhd2015_doctest", "runs",
		"2026-07-17-12-34-56-07-abc12def.jsonl",
	)
	if !strings.HasSuffix(resp.Path, wantSuffix) {
		t.Fatalf("path %q does not end with %q", resp.Path, wantSuffix)
	}
	if !strings.HasPrefix(resp.Path, req.CacheDir) {
		t.Fatalf("path %q not under cache %q", resp.Path, req.CacheDir)
	}
	base := filepath.Base(resp.Path)
	if base != "2026-07-17-12-34-56-07-abc12def.jsonl" {
		t.Fatalf("basename = %q", base)
	}
}
```
