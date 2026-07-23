## Expected

- `!e2e && heavy` matches only when heavy is present and e2e is absent.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "!e2e && heavy", []string{"heavy"})
	if err != nil || !ok {
		t.Fatalf("!e2e && heavy on {heavy}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && heavy", []string{"e2e", "heavy"})
	if err != nil || ok {
		t.Fatalf("!e2e && heavy on {e2e,heavy}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && heavy", nil)
	if err != nil || ok {
		t.Fatalf("!e2e && heavy on {}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && heavy", []string{"e2e"})
	if err != nil || ok {
		t.Fatalf("!e2e && heavy on {e2e}: ok=%v err=%v", ok, err)
	}
}
```
