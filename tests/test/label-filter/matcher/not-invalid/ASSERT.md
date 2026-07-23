## Expected

- Lone `!` and trailing `!` after a label are rejected.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	_, err = evalLabelExpr(t, "!", []string{"e2e"})
	if err == nil {
		t.Fatal("expected parse error for bare !")
	}
	_, err = evalLabelExpr(t, "e2e !", []string{"e2e"})
	if err == nil {
		t.Fatal("expected parse error for trailing !")
	}
	// keyword "not" is not supported
	_, err = evalLabelExpr(t, "not e2e", []string{"e2e"})
	if err == nil {
		t.Fatal("expected parse error for keyword not (unsupported)")
	}
}
```
