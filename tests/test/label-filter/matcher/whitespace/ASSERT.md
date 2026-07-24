## Expected

- Leading and trailing spaces do not break parsing.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, " slow && ui ", []string{"slow", "ui"})
	if err != nil || !ok {
		t.Fatalf("whitespace trim: ok=%v err=%v", ok, err)
	}
}
```