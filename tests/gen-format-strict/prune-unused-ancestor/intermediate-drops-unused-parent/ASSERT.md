## Expected

- Generate completes enough to write intermediate packages (`IntermediateSetupGo` non-empty preferred).
- Intermediate mid `setup.go` does **not** import the unused parent `feature` package
  (no path segment `/feature"` or alias import of the parent package for symbols never used).
- Suite may pass or fail independently; this leaf keys off **import prune**, not suite exit.
- droot import may remain when mid Setup still uses `droot.Request`.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}

	src := resp.IntermediateSetupGo
	if src == "" {
		// Fallback: find mid setup in snippets / walk-like fields.
		src = findMidSetup(resp.AllGoSnippets)
	}
	if src == "" {
		t.Fatalf("A3: intermediate mid/setup.go not found under GenDir=%s\nRunErr=%v\nsnippets head:\n%.800s",
			resp.GenDir, resp.RunErr, resp.AllGoSnippets)
	}

	// Parent intermediate package import path is testcase/.../feature (not .../feature/mid).
	// Reject imports whose path ends with /feature or is exactly .../feature.
	if hasUnusedFeatureParentImport(src) {
		t.Fatalf("A3: intermediate setup still imports unused parent feature package:\n%s", src)
	}
}

func findMidSetup(all string) string {
	// Snippets are prefixed with "// file: rel"
	parts := strings.Split(all, "// file: ")
	for _, p := range parts {
		if strings.Contains(p, "/mid/setup.go") || strings.HasPrefix(p, "mid/setup.go") ||
			strings.Contains(p, "feature/mid/setup.go") {
			// body after first newline
			if i := strings.Index(p, "\n"); i >= 0 {
				return p[i+1:]
			}
			return p
		}
	}
	return ""
}

func hasUnusedFeatureParentImport(src string) bool {
	// Look at import lines only.
	lines := strings.Split(src, "\n")
	inImport := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "import (") {
			inImport = true
			continue
		}
		if inImport && trim == ")" {
			inImport = false
			continue
		}
		if !inImport && !strings.HasPrefix(trim, "import ") {
			continue
		}
		// extract quoted path
		q1 := strings.Index(trim, `"`)
		if q1 < 0 {
			continue
		}
		q2 := strings.LastIndex(trim, `"`)
		if q2 <= q1 {
			continue
		}
		path := trim[q1+1 : q2]
		// Parent feature package (not feature/mid, not feature/mid/leaf).
		base := filepath.Base(path)
		if base == "feature" {
			return true
		}
		// e.g. testcase/feature when tree-scoped at root
		if strings.HasSuffix(path, "/feature") {
			return true
		}
	}
	return false
}
```
