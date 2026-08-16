## Expected

- Child leaf runs and returns Name `child`.

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
	if resp.Name != "child" {
		t.Fatalf("Name = %q, want child", resp.Name)
	}
}
```
