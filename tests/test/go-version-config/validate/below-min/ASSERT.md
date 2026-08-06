## Expected

- Validation fails with message containing `< 99.0`.

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
	if !strings.Contains(resp.ErrMsg, "< 99.0") {
		t.Fatalf("want below-min error, got %q", resp.ErrMsg)
	}
}
```
