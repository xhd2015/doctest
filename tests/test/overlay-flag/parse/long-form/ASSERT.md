## Expected

- Parse succeeds.
- `Response.Overlay` is absolute and equals the path after `--overlay`.

## Errors

- Classic TDD **RED** while `--overlay` is unrecognized.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected parse success for --overlay, got exit=%d err=%s", resp.ExitCode, errText(resp))
	}
	if !absLike(resp.Overlay) {
		t.Fatalf("Options.Overlay must be absolute non-empty, got %q", resp.Overlay)
	}
	var want string
	for i, a := range req.ParseArgs {
		if a == "--overlay" && i+1 < len(req.ParseArgs) {
			want = req.ParseArgs[i+1]
			break
		}
	}
	if want == "" {
		t.Fatal("SETUP must pass --overlay FILE")
	}
	wantAbs, aerr := filepath.Abs(want)
	if aerr != nil {
		t.Fatalf("abs want: %v", aerr)
	}
	if resp.Overlay != want && resp.Overlay != wantAbs &&
		filepath.Clean(resp.Overlay) != filepath.Clean(wantAbs) {
		t.Fatalf("Overlay=%q want %q (or abs %q)", resp.Overlay, want, wantAbs)
	}
}
```
