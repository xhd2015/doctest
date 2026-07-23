## Expected

- `!(e2e || flaky)` is true only when neither label is present.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "!(e2e || flaky)", nil)
	if err != nil || !ok {
		t.Fatalf("!(e2e || flaky) on {}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!(e2e || flaky)", []string{"heavy"})
	if err != nil || !ok {
		t.Fatalf("!(e2e || flaky) on {heavy}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!(e2e || flaky)", []string{"e2e"})
	if err != nil || ok {
		t.Fatalf("!(e2e || flaky) on {e2e}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!(e2e || flaky)", []string{"flaky"})
	if err != nil || ok {
		t.Fatalf("!(e2e || flaky) on {flaky}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "!(e2e || flaky)", []string{"e2e", "flaky"})
	if err != nil || ok {
		t.Fatalf("!(e2e || flaky) on {e2e,flaky}: ok=%v err=%v", ok, err)
	}
}
```
