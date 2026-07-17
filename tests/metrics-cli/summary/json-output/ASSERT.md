## Expected

- Exit code 0.
- stdout is valid JSON.

## Exit Code

- 0

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	if !isValidJSON(resp.Stdout) {
		t.Fatalf("summary --json stdout not valid JSON:\n%s", resp.Stdout)
	}
}
```
