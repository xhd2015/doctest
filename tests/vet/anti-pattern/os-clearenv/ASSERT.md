## Expected
- The command fails with non-zero exit.
- stderr contains `anti-pattern:` and names `os.Clearenv`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit")
	}
	if !strings.Contains(resp.Stderr, "anti-pattern:") {
		t.Fatalf("stderr missing anti-pattern: prefix:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "os.Clearenv") {
		t.Fatalf("stderr missing os.Clearenv:\n%s", resp.Stderr)
	}
}
```
