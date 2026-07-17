---
label: heavy
---

## Expected

- `a || b && c` matches `{a}` and `{b,c}` but not `{b}` alone.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	expr := "a || b && c"
	cases := []struct {
		labels []string
		want   bool
	}{
		{[]string{"a"}, true},
		{[]string{"b", "c"}, true},
		{[]string{"b"}, false},
	}
	for _, tc := range cases {
		ok, err := evalLabelExpr(t, expr, tc.labels)
		if err != nil {
			t.Fatalf("expr=%q labels=%v: %v", expr, tc.labels, err)
		}
		if ok != tc.want {
			t.Fatalf("expr=%q labels=%v: got %v want %v", expr, tc.labels, ok, tc.want)
		}
	}
}
```