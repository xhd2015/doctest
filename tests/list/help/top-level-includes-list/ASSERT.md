## Expected

- Exit code 0.
- stdout top-level usage lists `list` as its own command entry (not merely
  substring matches like `skill --list` or `--list-sessions`).

## Exit Code

- 0

```go
import (
	"regexp"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	out := resp.Stdout
	if !strings.Contains(out, "Usage: doctest") {
		t.Fatalf("missing top-level usage:\n%s", out)
	}
	// Command entry line: indentation + "list" as the command name token.
	// Avoids false positives from "skill --list" / "--list-sessions".
	cmdEntry := regexp.MustCompile(`(?m)^[ \t]+list\b`)
	if !cmdEntry.MatchString(out) {
		t.Fatalf("top-level help must list list as a command entry:\n%s", out)
	}
	for _, want := range []string{"test", "vet", "build"} {
		if !strings.Contains(out, want) {
			t.Fatalf("top-level help missing %q:\n%s", want, out)
		}
	}
}
```
