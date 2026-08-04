## Expected

- Second write succeeds.
- Gen `go.mod` exists again with `module testcase`.
- Cache files still present after rebuild.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after deleting gen go.mod failed: %v", err)
	}
	if !fileExists(filepath.Join(req.GenDir, "go.mod")) {
		t.Fatal("expected gen go.mod restored after integrity miss")
	}
	if !strings.Contains(resp.GoModContent, "module testcase") {
		t.Fatalf("restored go.mod missing module testcase:\n%s", resp.GoModContent)
	}
	requireSrcAndBridges(t, resp)
}
```
