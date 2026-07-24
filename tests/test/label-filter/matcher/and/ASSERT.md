## Expected

- AND matches only when all tokens are present.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "slow && ui", []string{"slow", "ui"})
	if err != nil || !ok {
		t.Fatalf("and match: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "slow && ui", []string{"slow"})
	if err != nil || ok {
		t.Fatalf("and partial: ok=%v err=%v", ok, err)
	}
}
```