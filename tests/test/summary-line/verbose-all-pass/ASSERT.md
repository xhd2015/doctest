## Expected

- Command succeeds.
- stderr shows the go test command includes `-v`.
- stdout ends with `PASS (1/1)` after verbose go-test output.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "go test -mod=mod -v") {
		t.Fatalf("expected stderr to contain 'go test -mod=mod -v', got:\n%s", resp.Stderr)
	}
	summary := findResultSummary(resp.Stdout)
	if stripANSI(strings.TrimSpace(summary)) != "PASS (1/1)" {
		t.Fatalf("expected PASS (1/1) summary after verbose output, got %q\nstdout:\n%s",
			summary, resp.Stdout)
	}
	summaryIsLastLine(t, resp.Stdout)
}
```