## Expected

- Embed script succeeds.
- Output is a single `.go` file containing `package assert`.
- Output does not contain `*_test.go` markers or test-only symbols.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("embed script failed: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("embed script error: %v", resp.Err)
	}
	src := string(resp.OutputBytes)
	if !strings.Contains(src, "package assert") {
		t.Fatalf("expected package assert in output, got:\n%s", src)
	}
	if strings.Contains(src, "_test.go") {
		t.Fatalf("output must not reference _test.go files")
	}
	if strings.Contains(src, "func Test") {
		t.Fatalf("output must not include test functions")
	}
}
```