---
label: heavy
---

## Expected

- No error (missing is not failure).
- `Hit == false` for a never-written key.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("GetPass missing must not error: %v", err)
	}
	if resp.Hit {
		t.Fatal("GetPass on missing key must be false")
	}
}
```
