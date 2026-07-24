---
label: heavy
---

## Expected

- Match succeeds.
- The inline `<any-of>` inside `<contains>` is interpreted as matcher syntax, not literal output.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
