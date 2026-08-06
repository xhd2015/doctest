## Expected

- OR matches when at least one token is present.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	for _, labels := range [][]string{{"slow"}, {"flaky"}} {
		ok, err := evalLabelExpr(t, "slow || flaky", labels)
		if err != nil || !ok {
			t.Fatalf("or match %v: ok=%v err=%v", labels, ok, err)
		}
	}
	ok, err := evalLabelExpr(t, "slow || flaky", []string{"fast"})
	if err != nil || ok {
		t.Fatalf("or miss: ok=%v err=%v", ok, err)
	}
}
```