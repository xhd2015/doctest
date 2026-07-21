---
label: heavy
explanation: nested doctest test compiles a sleep fixture and waits on go test -timeout=2s
---

## Expected

- Nested `go test` times out (sleep ≥ timeout).
- Process exits non-zero.
- Combined stdout+stderr contains a clear timeout signal findable without
  reading a full goroutine dump only, preferably:

  `Error: go test timed out after 2s`

  Accepted equivalents on the fail path:
  - `Error: go test timed out after …`
  - `test timed out after …`
  - `timed out after …`

## Errors

- Timeout is reported as a first-class fail-path message, not only buried
  under filtered JSON / progress dots.

## Exit Code

- Non-zero.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness run failed unexpectedly: %v\nstdout:\n%s\nstderr:\n%s",
			err, respStdout(resp), respStderr(resp))
	}
	if resp == nil {
		t.Fatal("expected non-nil Response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit when nested go test times out, got 0\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	combined := combinedOutput(resp)
	if !hasTimeoutSignal(combined) {
		t.Fatalf("expected clear timeout message (e.g. \"Error: go test timed out after 2s\") in stdout+stderr, got exit=%d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}

func respStdout(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stdout
}

func respStderr(resp *Response) string {
	if resp == nil {
		return ""
	}
	return resp.Stderr
}
```
