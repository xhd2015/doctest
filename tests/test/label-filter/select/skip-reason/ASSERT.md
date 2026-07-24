## Expected

- All five skipped entries have Reason "label filter".

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.SkippedReason) != 5 {
		t.Fatalf("reasons=%v", resp.SkippedReason)
	}
	for i, r := range resp.SkippedReason {
		if r != "label filter" {
			t.Fatalf("skipped[%d] reason=%q want %q (path=%s)", i, r, "label filter", resp.SkippedPaths[i])
		}
	}
}
```
