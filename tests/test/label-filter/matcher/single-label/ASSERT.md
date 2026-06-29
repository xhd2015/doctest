## Expected

- `slow` matches `{slow}` only.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	ok, err := evalLabelExpr(t, "slow", []string{"slow"})
	if err != nil || !ok {
		t.Fatalf("slow on {slow}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "slow", nil)
	if err != nil || ok {
		t.Fatalf("slow on {}: ok=%v err=%v", ok, err)
	}
	ok, err = evalLabelExpr(t, "slow", []string{"fast"})
	if err != nil || ok {
		t.Fatalf("slow on {fast}: ok=%v err=%v", ok, err)
	}
}
```