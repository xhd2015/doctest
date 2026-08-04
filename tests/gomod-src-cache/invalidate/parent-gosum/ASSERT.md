## Expected

- Second write succeeds after parent go.sum appears.
- Gen `go.sum` exists and reflects new content.
- `doctest.tidy-done` is absent (gen sum wrote).
- Fingerprint file still present after rebuild.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after parent go.sum change failed: %v", err)
	}
	sumPath := filepath.Join(req.GenDir, "go.sum")
	if !fileExists(sumPath) {
		t.Fatal("expected gen go.sum after parent go.sum write")
	}
	sumBody := readFileOrEmpty(sumPath)
	if !strings.Contains(sumBody, "example.com/dep") {
		t.Fatalf("expected gen go.sum to copy parent sum, got:\n%s", sumBody)
	}
	if resp.TidyDoneExists {
		t.Fatal("expected doctest.tidy-done dropped when gen go.sum wrote")
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src after sum invalidation rebuild")
	}
}
```
