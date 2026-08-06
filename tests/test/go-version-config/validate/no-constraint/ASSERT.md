## Expected

- Validation succeeds with empty ErrMsg.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ErrMsg != "" {
		t.Fatalf("want no validate error, got %q", resp.ErrMsg)
	}
	if !resp.Loaded {
		t.Fatal("want loaded empty config")
	}
}
```
