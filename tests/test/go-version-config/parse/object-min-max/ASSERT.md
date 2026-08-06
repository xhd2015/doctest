## Expected

- Load succeeds with Go.Min `1.18` and Go.Max `1.20`.

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
		t.Fatalf("load error: %s", resp.ErrMsg)
	}
	if !resp.Loaded || resp.Min != "1.18" || resp.Max != "1.20" {
		t.Fatalf("want min=1.18 max=1.20 loaded, got loaded=%v min=%q max=%q", resp.Loaded, resp.Min, resp.Max)
	}
}
```
