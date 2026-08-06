---
label: e2e, heavy
---

## Expected (TDD RED under external gen until product-internal strategy applies)

- Subject `doctest test` that **imports product** `example.com/app/internal/…`
  from a **runner** module (not in-module compile) **succeeds**.
- Today this **fails** with `use of internal package example.com/app/internal/…`
  under external unified gen — outer leaf is RED for that configuration.
- Success path for product internal remains `tests/build/in-module/` (WorkDir = app).

## Exit Code

- Subject doctest exit **0** (desired for this recipe once fixed, or switch to in-module).

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
		hint := ""
		if containsInternalDenied(resp.Combined) && strings.Contains(resp.Combined, "example.com/app/internal") {
			hint = " (kind B: product internal via external gen)"
		}
		t.Fatalf("kind B RED: want subject doctest exit 0, got %d%s\n%s", resp.ExitCode, hint, resp.Combined)
	}
}
```
