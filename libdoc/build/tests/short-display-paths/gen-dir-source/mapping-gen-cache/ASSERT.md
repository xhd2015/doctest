---
label: heavy
---

## Expected

- `build.Test` succeeds (`resp.TestErr` is nil).
- Process cwd is unchanged (`CwdBefore == CwdAfter`) — no process Chdir in Run.
- `resp.HeaderLine` equals `doctest: ` + `pathfmt.Short(resp.TestRoot)` + ` (1 tests)`
  (sandbox is usually outside process cwd → absolute; still names `tests/feature`).
- `resp.CdLine` contains `~/` and `mapping-gen` (home cache shortening; valid without Chdir).
- `resp.CdLine` does **not** contain the raw absolute home path.

## Exit Code

- `err` is nil.

```go
import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.TestErr != nil {
		t.Fatalf("expected build.Test to succeed, got: %v\nstderr:\n%s", resp.TestErr, resp.Stderr)
	}
	assertNoProcessChdir(t, resp)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	wantHeader := fmt.Sprintf("doctest: %s (1 tests)", pathfmt.Short(resp.TestRoot))
	if resp.HeaderLine != wantHeader {
		t.Fatalf("header must be pathfmt.Short(testRoot) form\nwant %q\ngot  %q\nstderr:\n%s",
			wantHeader, resp.HeaderLine, resp.Stderr)
	}
	// Still identify the fixture tree segment (compact relative or abs suffix).
	if !strings.Contains(resp.HeaderLine, "tests/feature") &&
		!strings.Contains(resp.HeaderLine, "tests"+string(os.PathSeparator)+"feature") {
		t.Fatalf("header must name tests/feature, got %q", resp.HeaderLine)
	}

	if resp.CdLine == "" {
		t.Fatalf("missing cd line in stderr:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.CdLine, home) {
		t.Fatalf("cd line must not contain raw home %q: %s", home, resp.CdLine)
	}
	if !strings.Contains(resp.CdLine, "mapping-gen") {
		t.Fatalf("cd line must contain mapping-gen: %s", resp.CdLine)
	}
	if !strings.Contains(resp.CdLine, "~/") {
		t.Fatalf("cd line must use ~/ for cache path: %s", resp.CdLine)
	}
}
```
