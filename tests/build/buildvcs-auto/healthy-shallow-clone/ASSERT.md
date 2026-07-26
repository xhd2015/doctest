## Expected

- Shallow checkout does **not** make `-buildvcs=auto` fail.
- `go build -buildvcs=auto` exits 0 and stamps VCS metadata (tip commit).

## Exit Code

- 0

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("shallow clone: expected exit 0 under -buildvcs=auto (shallow is not the fail axis), got %d\n%s",
			resp.ExitCode, resp.Output)
	}
	assertVCSStamped(t, resp.OutBin)
}
```
