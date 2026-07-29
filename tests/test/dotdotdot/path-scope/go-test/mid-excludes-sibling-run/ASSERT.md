---
label: heavy
---

## Expected

Mid path go test executes `MARKER:MID_LEAF` and never `MARKER:SIBLING_LEAF`.

```go
import (
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
	if !strings.Contains(out, "MARKER:MID_LEAF") {
		t.Fatalf("want mid leaf executed\n%s", out)
	}
	if strings.Contains(out, "MARKER:SIBLING_LEAF") {
		t.Fatalf("sibling must not run under mid scope\n%s", out)
	}
}
```
