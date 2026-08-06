## Expected

- `!e2e` is true when the leaf does not carry `e2e` (including unlabeled).
- `!e2e` is false when `e2e` is present.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "!e2e", nil)
	if err != nil || !ok {
		t.Fatalf("!e2e on {}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e", []string{"flaky"})
	if err != nil || !ok {
		t.Fatalf("!e2e on {flaky}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e", []string{"e2e"})
	if err != nil || ok {
		t.Fatalf("!e2e on {e2e}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e", []string{"e2e", "flaky"})
	if err != nil || ok {
		t.Fatalf("!e2e on {e2e,flaky}: ok=%v err=%v", ok, err)
	}
	// double negation
	ok, err = evalLabelExpr(t, "!!e2e", []string{"e2e"})
	if err != nil || !ok {
		t.Fatalf("!!e2e on {e2e}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!!e2e", nil)
	if err != nil || ok {
		t.Fatalf("!!e2e on {}: ok=%v err=%v", ok, err)
	}
}
```
