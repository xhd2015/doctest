## Expected

- Exit code 0.
- stderr contains `no tests changed`.

## Side Effects

- No tests are executed.

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "no tests changed") {
		t.Fatalf("stderr missing 'no tests changed':\n%s", resp.Stderr)
	}
}
```