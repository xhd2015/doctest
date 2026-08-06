## Expected

- Non-zero exit.
- Error indicates unknown/unrecognized flag or that overlay is not accepted on build.
- Must **not** succeed as a normal build.

## Side Effects

- None required. Regression: stays failing after `test` accepts `-overlay`.

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
		t.Fatal("doctest build must still reject -overlay")
	}
	msg := strings.ToLower(errText(resp) + resp.Stdout)
	// Accept current lessflags "unrecognized" or a future explicit "not supported".
	if !strings.Contains(msg, "overlay") &&
		!strings.Contains(msg, "unrecognized") &&
		!strings.Contains(msg, "unknown") {
		t.Fatalf("expected reject message about overlay/unknown flag, got:\n%s", errText(resp)+resp.Stdout)
	}
}
```
