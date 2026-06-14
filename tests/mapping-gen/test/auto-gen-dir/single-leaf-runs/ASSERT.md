## Expected
- Command succeeds (exit code 0).
- The test case runs and passes.
- Gen dir path is printed to stderr (begins with `→ `).
- Gen dir is under a `mapping-gen/` cache location.

## Exit Code
- Exit code 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	genDirPath := parseGenDir(resp.Stderr)
	if genDirPath == "" {
		t.Fatalf("expected gen dir path in stderr (→ <path>):\n%s", resp.Stderr)
	}

	if !strings.Contains(genDirPath, "mapping-gen") {
		t.Fatalf("expected mapping-gen in gen dir path, got: %s", genDirPath)
	}

	if !strings.Contains(resp.Stdout, "PASS") {
		t.Fatalf("expected PASS in stdout:\n%s", resp.Stdout)
	}
}
```
