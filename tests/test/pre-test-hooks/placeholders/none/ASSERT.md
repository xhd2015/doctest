## Expected

- The command runs from the project root unchanged.
- Neither overlay artifact is allocated and no Go flag is produced.

```go
import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if !reflect.DeepEqual(resp.Calls, [][]string{{"tool", "prepare"}}) { t.Fatalf("calls=%#v", resp.Calls) }
	if !reflect.DeepEqual(resp.WorkDirs, []string{filepath.Join(d.DOCTEST_ROOT, "project")}) { t.Fatalf("workdirs=%#v", resp.WorkDirs) }
	if resp.OverlayDir != "" || resp.OverlayFile != "" || len(resp.GoFlags) != 0 { t.Fatalf("unexpected overlay state: %#v", resp) }
}
```
