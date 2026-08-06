## Expected

- `build.Test` succeeds.
- Process cwd is unchanged (`CwdBefore == CwdAfter`) — no process Chdir in Run.
- `resp.ResolvedGenDir` is absolute under `resp.ProjRoot` and ends with `_gen`.
- `resp.HeaderLine` equals `doctest: ` + `pathfmt.Short(resp.TestRoot)` + ` (1 tests)`.
- `resp.CdLine` applies `pathfmt.Short` to the explicit gen/run dir and still
  names the `_gen` segment (without Chdir, Short often yields the absolute
  sandbox path — that is correct; compact bare `_gen` required process cwd =
  project root, which is Parallel-unsafe).

## Exit Code

- `err` is nil.

```go
import (
	"fmt"
	"path/filepath"
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

	if resp.ResolvedGenDir == "" {
		t.Fatal("ResolvedGenDir empty for explicit gen-dir leaf")
	}
	if !filepath.IsAbs(resp.ResolvedGenDir) {
		t.Fatalf("ResolvedGenDir must be absolute (no Chdir-relative gen): %q", resp.ResolvedGenDir)
	}
	wantGen := filepath.Join(resp.ProjRoot, "_gen")
	if filepath.Clean(resp.ResolvedGenDir) != filepath.Clean(wantGen) {
		t.Fatalf("ResolvedGenDir want %q got %q", wantGen, resp.ResolvedGenDir)
	}

	wantHeader := fmt.Sprintf("doctest: %s (1 tests)", pathfmt.Short(resp.TestRoot))
	if resp.HeaderLine != wantHeader {
		t.Fatalf("header must be pathfmt.Short(testRoot) form\nwant %q\ngot  %q\nstderr:\n%s",
			wantHeader, resp.HeaderLine, resp.Stderr)
	}

	if resp.CdLine == "" {
		t.Fatalf("missing cd line in stderr:\n%s", resp.Stderr)
	}
	// DisplayPath short form of the explicit gen (or a run dir under it).
	// Without process Chdir, Short of a temp sandbox gen is typically absolute
	// (still valid DisplayPath); compact bare `_gen` is not required.
	shortGen := pathfmt.Short(resp.ResolvedGenDir)
	if !strings.Contains(resp.CdLine, shortGen) && !strings.Contains(resp.CdLine, "_gen") {
		t.Fatalf("cd line must use Short(explicit gen) or contain _gen segment\nshortGen=%q\ncd=%s\nstderr:\n%s",
			shortGen, resp.CdLine, resp.Stderr)
	}
}
```
