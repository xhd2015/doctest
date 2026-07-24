## Expected

- Trailing `&&` is rejected.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	_, err = evalLabelExpr(t, "slow &&", []string{"slow"})
	if err == nil {
		t.Fatal("expected parse error for trailing &&")
	}
}
```