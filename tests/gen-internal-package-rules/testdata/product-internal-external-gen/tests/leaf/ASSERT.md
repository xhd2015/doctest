## Expected

- Message from product internal greet.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if resp == nil || resp.Message != "hello-from-app-internal" {
		t.Fatalf("Message = %v", resp)
	}
}
```
