---
label: e2e, heavy
---

## Expected (TDD RED until kind A shim)

- Subject `doctest test` on a tree with leaf path `http/internal/…` **succeeds**
  (gen suite packages).
- Today this **fails** at go test with `use of internal package` on the **gen**
  path — outer leaf is RED until shim lands.

## Exit Code

- Subject doctest exit **0** (desired).

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run harness error: %v\n%s", err, resp.Combined)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		// Help distinguish kind A packaging failure while RED.
		msg := resp.Combined
		hint := ""
		if containsInternalDenied(msg) {
			hint = " (kind A packaging: gen path under …/internal/…)"
			if strings.Contains(msg, "example.com/app/internal") {
				hint = " (unexpected product internal — looks like kind B)"
			}
		}
		t.Fatalf("kind A RED: want subject doctest exit 0, got %d%s\n%s", resp.ExitCode, hint, msg)
	}
}
```
