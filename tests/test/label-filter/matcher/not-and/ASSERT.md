## Expected

- `!e2e && flaky` matches only when flaky is present and e2e is absent.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "!e2e && flaky", []string{"flaky"})
	if err != nil || !ok {
		t.Fatalf("!e2e && flaky on {flaky}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && flaky", []string{"e2e", "flaky"})
	if err != nil || ok {
		t.Fatalf("!e2e && flaky on {e2e,flaky}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && flaky", nil)
	if err != nil || ok {
		t.Fatalf("!e2e && flaky on {}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!e2e && flaky", []string{"e2e"})
	if err != nil || ok {
		t.Fatalf("!e2e && flaky on {e2e}: ok=%v err=%v", ok, err)
	}
}
```
