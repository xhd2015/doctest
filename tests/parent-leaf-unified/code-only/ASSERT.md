## Expected

- Parent leaf runs and returns Name `code-only`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "code-only" {
		t.Fatalf("Name = %q, want code-only", resp.Name)
	}
}
```
