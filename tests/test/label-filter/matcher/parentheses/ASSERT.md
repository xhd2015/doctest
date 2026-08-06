## Expected

- Parentheses require `ui` and at least one of `slow` or `flaky`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	expr := "(slow || flaky) && ui"
	ok, err := evalLabelExpr(t, expr, []string{"slow", "ui"})
	if err != nil || !ok {
		t.Fatalf("slow+ui: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, expr, []string{"flaky", "ui"})
	if err != nil || !ok {
		t.Fatalf("flaky+ui: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, expr, []string{"slow"})
	if err != nil || ok {
		t.Fatalf("slow only: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, expr, []string{"ui"})
	if err != nil || ok {
		t.Fatalf("ui only: ok=%v err=%v", ok, err)
	}
}
```