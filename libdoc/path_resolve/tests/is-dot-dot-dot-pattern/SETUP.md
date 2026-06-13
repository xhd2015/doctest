## Preconditions
- The `IsDotDotDotPattern` function is defined in `path_resolve`.
- Input is a string argument (e.g. `./foo/...`, `foo/bar`, `...`).

## Steps
1. Call `path_resolve.IsDotDotDotPattern(req.Input)`.
2. Return the result in `resp.BoolResult`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/path_resolve"
)

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	result := path_resolve.IsDotDotDotPattern(req.Input)
	return &Response{BoolResult: result}, nil
}
```
