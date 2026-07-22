## Expected

- Labels = [ui-automation]; Explanation set.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse err: %s", resp.ParseErr)
	}
	if len(resp.Labels) != 1 || resp.Labels[0] != "ui-automation" {
		t.Fatalf("labels=%v", resp.Labels)
	}
	if resp.Explanation != "heavy ui test" {
		t.Fatalf("explanation=%q", resp.Explanation)
	}
}
```
