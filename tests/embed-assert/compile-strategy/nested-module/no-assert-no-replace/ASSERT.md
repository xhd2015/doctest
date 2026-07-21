---
label: heavy
---

## Expected

- `doctest test` exits 0.
- Nested `go.mod` under outside gen-dir **contains** `replace github.com/xhd2015/doctest/assert =>` pointing at assert-mod cache (always-on for external modules).

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected no-author-assert nested-module test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNestedGoMod(t, req.OutsideGenDir)
	assertAssertReplaceInGoMod(t, req.OutsideGenDir+"/go.mod")
}
```
