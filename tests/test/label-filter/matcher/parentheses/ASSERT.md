---
label: heavy
---

## Expected

- Parentheses require `ui` and at least one of `slow` or `heavy`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	expr := "(slow || heavy) && ui"
	ok, err := evalLabelExpr(t, expr, []string{"slow", "ui"})
	if err != nil || !ok {
		t.Fatalf("slow+ui: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, expr, []string{"heavy", "ui"})
	if err != nil || !ok {
		t.Fatalf("heavy+ui: ok=%v err=%v", ok, err)
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