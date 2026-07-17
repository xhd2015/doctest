## Expected

- Message is non-empty.
- Contains `WARNING:`.
- Contains `skill:doctest-review-perf`.
- Contains `review-perf --show` (skill invocation hint).
- Contains `3 minutes`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := resp.Message
	if msg == "" {
		t.Fatal("FormatDefaultSuiteSlowWarning returned empty string")
	}
	needles := []string{
		"WARNING:",
		"skill:doctest-review-perf",
		"review-perf --show",
		"3 minutes",
	}
	for _, n := range needles {
		if !strings.Contains(msg, n) {
			t.Fatalf("message missing %q\nfull:\n%s", n, msg)
		}
	}
}
```
