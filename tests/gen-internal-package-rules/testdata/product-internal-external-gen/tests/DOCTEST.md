# expose subject — product internal import under external gen

Leaf imports `example.com/app/internal/greet`. Outer harness runs this tree
from a **runner** module so doctest does not enter product internal-compile.

```go
import (
	"testing"

	"example.com/app/internal/greet"
	"github.com/xhd2015/doctest/session"
)

type Request struct{}

type Response struct {
	Message string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	_ = req
	return &Response{Message: greet.Hello()}, nil
}
```
