## Expected

- The command succeeds and prints skill-level help instead of a version.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr: %s", resp.ExitCode, resp.Stderr)
	}
	for _, want := range []string{"Usage: doctest skill --list", "--version", "--install"} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	if resp.Stdout == "0.1.0\n" {
		t.Fatalf("stdout must be skill help, got version only")
	}
}
```
