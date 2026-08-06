## Expected

`./tree/mid/...`: mid + nested under prefix; **no** sibling.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, out)
	}
	if !strings.Contains(out, "MARKER:MID_LEAF") || !strings.Contains(out, "MARKER:NESTED_LEAF") {
		t.Fatalf("want mid+nested under prefix\n%s", out)
	}
	if strings.Contains(out, "MARKER:SIBLING_LEAF") {
		t.Fatalf("sibling must not run\n%s", out)
	}
}
```
