## Expected

- Hooks execute in declaration order from the same project root.
- Both commands receive byte-for-byte identical substituted directory and file paths.
- An untouched pre-created file still contributes no Go overlay flag.

```go
import (
	"reflect"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if len(resp.Calls) != 2 || resp.Calls[0][0] != "first" || resp.Calls[1][0] != "second" { t.Fatalf("call order=%#v", resp.Calls) }
	if !reflect.DeepEqual(resp.WorkDirs, []string{resp.WorkDirs[0], resp.WorkDirs[0]}) { t.Fatalf("workdirs not shared=%#v", resp.WorkDirs) }
	for _, at := range []int{2, 4} {
		if resp.Calls[0][at] != resp.Calls[1][at] { t.Fatalf("placeholder arg %d differs: %#v", at, resp.Calls) }
	}
	if resp.Calls[0][2] != resp.OverlayDir || resp.Calls[0][4] != resp.OverlayFile || len(resp.GoFlags) != 0 { t.Fatalf("unexpected shared state: %#v", resp) }
}
```
