## Expected

- Labels empty; Explanation set; no parse error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse err: %s", resp.ParseErr)
	}
	if len(resp.Labels) != 0 {
		t.Fatalf("labels=%v want empty", resp.Labels)
	}
	if resp.Explanation != "documentation note only" {
		t.Fatalf("explanation=%q", resp.Explanation)
	}
}
```
