---
explanation: nested go test with memprofile
---

## Expected
- Exit code 0.
- stderr contains `-memprofile=` with the absolute path supplied on the CLI.

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
	var want string
	for i, a := range req.Args {
		if a == "-memprofile" && i+1 < len(req.Args) {
			want = req.Args[i+1]
			break
		}
	}
	if want == "" {
		t.Fatal("could not determine expected memprofile path from req.Args")
	}
	if !strings.Contains(resp.Stderr, "-memprofile="+want) &&
		!strings.Contains(resp.Stderr, "-memprofile "+want) {
		t.Fatalf("expected stderr to contain -memprofile=%s, got:\n%s", want, resp.Stderr)
	}
}
```
