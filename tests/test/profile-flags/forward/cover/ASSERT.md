---
explanation: nested go test with -cover and -coverprofile
---

## Expected
- Exit code 0.
- stderr contains `-cover` as a go test flag.
- stderr contains abs-resolved `-coverprofile=` under WorkDir.

```go
import (
	"path/filepath"
	"regexp"
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
	// -cover as its own flag (not only as substring of -coverprofile).
	// Match word-ish: space or start, then -cover, then space/end/= not "profile".
	coverFlag := regexp.MustCompile(`(?:^|[\s])-cover(?:[\s]|$)`)
	if !coverFlag.MatchString(resp.Stderr) {
		// Also accept -cover=true style if emitted that way.
		if !strings.Contains(resp.Stderr, "-cover=true") {
			t.Fatalf("expected -cover on go command line, stderr:\n%s", resp.Stderr)
		}
	}
	wantAbs, err := filepath.Abs(filepath.Join(req.WorkDir, "cov", "cover.out"))
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if !strings.Contains(resp.Stderr, "-coverprofile="+wantAbs) &&
		!strings.Contains(resp.Stderr, "-coverprofile "+wantAbs) {
		t.Fatalf("expected -coverprofile=%s, stderr:\n%s", wantAbs, resp.Stderr)
	}
}
```
