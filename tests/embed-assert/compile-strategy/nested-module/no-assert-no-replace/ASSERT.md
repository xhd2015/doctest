## Expected

- `doctest test` exits 0.
- Nested `go.mod` has parent module replace but no assert replace.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected no-assert nested-module test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNestedGoMod(t, req.OutsideGenDir)
	assertNoAssertReplaceInGoMod(t, req.OutsideGenDir+"/go.mod")
}
```