## Expected

- The command succeeds.
- stdout is the raw declared version plus a newline.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr: %s", resp.ExitCode, resp.Stderr)
	}
	if resp.Stdout != "0.1.0\n" {
		t.Fatalf("stdout = %q, want %q", resp.Stdout, "0.1.0\\n")
	}
}
```
