## Expected

- Non-zero exit / parse error.
- Error text mentions `overlay` (and preferably missing argument / value).

## Errors

- Classic TDD: before the flag exists, may be "unrecognized flag" (still non-zero).
  Once the flag exists, missing value must still fail (arity).

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for -overlay without value")
	}
	msg := strings.ToLower(errText(resp))
	if !strings.Contains(msg, "overlay") {
		t.Fatalf("error must mention overlay, got:\n%s", errText(resp))
	}
}
```
