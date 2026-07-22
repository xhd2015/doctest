---
label: heavy
---

## Expected

- Suite generate + `go test` succeed (`resp.RunErr` empty).
- **Do not** require that generated sources equal `gofmt` output.
- Success is compile/run only (Phase A: format.Source not required).

```go
import (
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
	if resp.RunErr != "" {
		t.Fatalf("A6: expected compile/run success with explicit imports, RunErr=%v\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr)
	}
	// Explicitly do NOT fail when sources are not gofmt-pretty.
	// (Optional log only — never t.Fatal on gofmt mismatch.)
	if resp.LeafGo != "" && !looksGofmtEqual(resp.LeafGo) {
		t.Logf("A6 note: leaf Go is not gofmt-identical (allowed under Phase A)")
	}
}
```
