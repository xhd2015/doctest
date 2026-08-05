---
label: heavy
explanation: nested go test forwards covermode/coverpkg/short/failfast/parallel/shuffle/tags/gcflags/ldflags/race
---

## Expected
- Exit code 0 (fixture passes; tags/gcflags/ldflags are harmless here).
- Verbose stderr go command line includes each forwarded flag form.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	// Prefer stderr (verbose go line); fall back to combined.
	out := resp.Stderr + "\n" + resp.Stdout
	want := []string{
		"-covermode=atomic",
		"-coverpkg=example.com/mod/...",
		"-cover",
		"-short",
		"-failfast",
		"-parallel=2",
		"-shuffle=on",
		"-tags=integration",
		"-gcflags=all=-N",
		"-ldflags=-s",
		"-race",
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			// Space form is also acceptable for some flags.
			space := strings.Replace(w, "=", " ", 1)
			if space != w && strings.Contains(out, space) {
				continue
			}
			t.Fatalf("expected %q on go command line, out:\n%s", w, out)
		}
	}
}
```
