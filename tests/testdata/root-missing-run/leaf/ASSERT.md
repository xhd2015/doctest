---
label: heavy
---

## Expected
- This leaf would check that Run produces expected output if it existed.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    _ = req
    _ = resp
    _ = err
}
```
