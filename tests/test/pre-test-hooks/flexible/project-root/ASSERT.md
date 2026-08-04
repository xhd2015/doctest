## Expected

- `$PROJECT_ROOT` in a command arg expands to the absolute project root path.
- Unrelated `$OTHER` tokens are not expanded from the process environment.
- No overlay directory, file, or Go overlay flag is allocated.

```go
import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if resp.OverlayDir != "" || resp.OverlayFile != "" || len(resp.GoFlags) != 0 {
		t.Fatalf("unexpected overlay state: %#v", resp)
	}
	root, absErr := filepath.Abs(filepath.Join(d.DOCTEST_ROOT, "project"))
	if absErr != nil {
		t.Fatalf("abs project root: %v", absErr)
	}
	wantConfig := "--config=" + root + "/cfg"
	want := [][]string{{"tool", wantConfig, "--literal=$OTHER"}}
	if !reflect.DeepEqual(resp.Calls, want) {
		t.Fatalf("calls=%#v want %#v", resp.Calls, want)
	}
	if !reflect.DeepEqual(resp.WorkDirs, []string{filepath.Join(d.DOCTEST_ROOT, "project")}) {
		t.Fatalf("workdirs=%#v", resp.WorkDirs)
	}
}
```
