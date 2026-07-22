---
label: heavy
---

## Expected

- `Run` returns no harness error (compile failure is in `resp.RunErr`).
- Suite/generate path **fails** (`resp.RunErr` non-empty).
- Failure is consistent with missing `time` (stderr/run err mention `time` / undefined / undeclared), **or** generated Go that references `time.` does not contain `import "time"` / `"time"`.
- Engine must **not** inject `"time"` solely because the body used `time.Second`.

```go
import (
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
	if resp.RunErr == "" {
		t.Fatalf("expected generate/suite to fail without import \"time\", but RunErr empty\nstdout:\n%s\nstderr:\n%s\nleaf:\n%s",
			resp.Stdout, resp.Stderr, resp.LeafGo)
	}

	// Strong signal: generated sources that use time. must not gain "time" via auto-add.
	combined := resp.LeafGo + "\n" + resp.AllGoSnippets + "\n" + resp.RunErr + "\n" + resp.Stderr
	usesTimeSel := strings.Contains(resp.LeafGo, "time.") || strings.Contains(resp.AllGoSnippets, "time.Second")
	hasTimeImp := containsImportTime(resp.LeafGo) || containsImportTime(resp.AllGoSnippets)

	if usesTimeSel && hasTimeImp {
		t.Fatalf("A1: engine auto-added import \"time\" for user time.Second usage; want missing import + compile fail\nleaf:\n%s\nrunErr: %s",
			resp.LeafGo, resp.RunErr)
	}

	// Soft check: error path mentions time when available.
	low := strings.ToLower(combined)
	if !strings.Contains(low, "time") {
		t.Logf("note: failure output did not mention time; RunErr=%q stderr=%q", resp.RunErr, resp.Stderr)
	}
}
```
