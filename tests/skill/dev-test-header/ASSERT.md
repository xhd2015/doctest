## Expected

- The command succeeds.
- The header contains the standard nested `metadata.version` value.

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
	for _, want := range []string{"---\n", "name: doctest-dev-test", "metadata:\n  version: \"0.1.0\""} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	if strings.Contains(resp.Stdout, "# Dev-test") {
		t.Fatalf("header output includes skill body:\n%s", resp.Stdout)
	}
}
```
