## Expected
- Command succeeds (exit code 0).
- `go build` compiles successfully.
- Gen dir is printed to stderr.

## Exit Code
- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("build failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	genDirPath := parseGenDir(resp.Stderr)
	if genDirPath == "" {
		t.Fatalf("expected gen dir path in stderr (→ <path>):\n%s", resp.Stderr)
	}

	if !strings.Contains(resp.Stderr, "go build") {
		t.Fatalf("expected 'go build' in stderr:\n%s", resp.Stderr)
	}
}
```
