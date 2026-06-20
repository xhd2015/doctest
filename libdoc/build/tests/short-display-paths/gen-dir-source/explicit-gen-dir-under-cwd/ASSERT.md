## Expected

- `build.Test` succeeds.
- `resp.ArrowLine` is `→ ./_gen` (platform-native separator after `./`).
- `resp.ArrowLine` does **not** contain the temp absolute project path.

## Exit Code

- `err` is nil.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected build.Test to succeed, got: %v\nstderr:\n%s", resp.TestErr, resp.Stderr)
	}
	if resp.ArrowLine == "" {
		t.Fatalf("missing arrow line in stderr:\n%s", resp.Stderr)
	}
	wantSuffix := "." + string(filepath.Separator) + "_gen"
	arrowPath := strings.TrimSpace(strings.TrimPrefix(resp.ArrowLine, "→ "))
	if arrowPath != wantSuffix {
		t.Fatalf("expected arrow path %q, got %q (full line: %s)", wantSuffix, arrowPath, resp.ArrowLine)
	}
	if strings.Contains(resp.ArrowLine, "/var/") || strings.Contains(resp.ArrowLine, "\\Temp\\") {
		t.Fatalf("explicit gen dir must not show temp absolute path: %s", resp.ArrowLine)
	}
}
```