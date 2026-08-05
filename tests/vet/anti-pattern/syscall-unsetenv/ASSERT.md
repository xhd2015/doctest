## Expected
- The command fails with non-zero exit.
- stderr contains `anti-pattern:` and names `syscall.Unsetenv`.

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
	if !strings.Contains(resp.Stderr, "syscall.Unsetenv") {
		t.Fatalf("stderr missing syscall.Unsetenv:\n%s", resp.Stderr)
	}
}
```
