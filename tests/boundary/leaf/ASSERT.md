---
label: heavy
---

## Expected
- Run returns a non-nil error because the stub returns `fmt.Errorf("stub: not implemented")`.

## Errors
- `err` must not be nil.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err == nil {
        t.Fatal("expected error from stub Run, got nil")
    }
}
```
