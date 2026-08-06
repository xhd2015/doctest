## Expected

- Validation fails with message containing `> 1.0`.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if !strings.Contains(resp.ErrMsg, "> 1.0") {
		t.Fatalf("want above-max error, got %q", resp.ErrMsg)
	}
}
```
