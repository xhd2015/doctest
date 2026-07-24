## Expected
- `FindDotDotDotDirs(".")` returns both the parent and nested `DOCTEST.md` roots.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.DirsResult) != 2 {
		t.Fatalf("expected 2 doctest roots, got %d: %v", len(resp.DirsResult), resp.DirsResult)
	}

	var hasParent, hasNested bool
	for _, dir := range resp.DirsResult {
		base := filepath.Base(dir)
		switch base {
		case filepath.Base(req.BasePath):
			hasParent = true
		case "mapping-gen":
			hasNested = true
		default:
			t.Fatalf("unexpected doctest root %q in %v", base, resp.DirsResult)
		}
	}
	if !hasParent || !hasNested {
		t.Fatalf("expected parent and mapping-gen roots, got %v", resp.DirsResult)
	}
	if strings.Contains(resp.DirsResult[0], "mapping-gen") && strings.Contains(resp.DirsResult[1], "mapping-gen") {
		t.Fatalf("expected only one nested root, got %v", resp.DirsResult)
	}
}
```