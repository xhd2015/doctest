## Expected

- Exit code 0.
- stdout documents `doctest list` usage.
- Mentions patterns (`./...` or path patterns).
- Explains L2:L3 (e2e vs non-e2e leaves).
- Documents `--color` and `--no-color`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	out := resp.Stdout
	for _, want := range []string{
		"Usage:",
		"list",
		"./...",
		"L2",
		"L3",
		"--color",
		"--no-color",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("list --help missing %q:\n%s", want, out)
		}
	}
	if !strings.Contains(strings.ToLower(out), "e2e") {
		t.Fatalf("list --help should mention e2e (L3 identity):\n%s", out)
	}
}
```
