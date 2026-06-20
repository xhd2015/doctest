## Expected

- `build.Test` succeeds (`resp.TestErr` is nil).
- `resp.ArrowLine` contains `mapping-gen` and starts with `→ ~` (home-shortened).
- `resp.ArrowLine` does **not** contain the raw absolute home path.
- `resp.HeaderLine` starts with `doctest: ./` (test tree under cwd).
- `resp.CdLine` contains `~/` and `mapping-gen` (gen run dir under home cache).

## Exit Code

- `err` is nil.

```go
import (
	"os"
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
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if resp.ArrowLine == "" {
		t.Fatalf("missing arrow line in stderr:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.ArrowLine, home) {
		t.Fatalf("arrow line must not contain raw home %q: %s", home, resp.ArrowLine)
	}
	if !strings.Contains(resp.ArrowLine, "mapping-gen") {
		t.Fatalf("arrow line must contain mapping-gen: %s", resp.ArrowLine)
	}
	arrowPath := strings.TrimPrefix(resp.ArrowLine, "→ ")
	if !strings.HasPrefix(arrowPath, "~") {
		t.Fatalf("arrow path must start with ~, got %q", arrowPath)
	}
	if resp.HeaderLine == "" {
		t.Fatalf("missing doctest line in stderr:\n%s", resp.Stderr)
	}
	if !strings.HasPrefix(resp.HeaderLine, "doctest: ./") {
		t.Fatalf("doctest line must use ./ prefix, got %q", resp.HeaderLine)
	}
	if resp.CdLine == "" {
		t.Fatalf("missing cd line in stderr:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.CdLine, "mapping-gen") {
		t.Fatalf("cd line must contain mapping-gen: %s", resp.CdLine)
	}
	if !strings.Contains(resp.CdLine, "~/") {
		t.Fatalf("cd line must use ~/ for cache path: %s", resp.CdLine)
	}
}
```