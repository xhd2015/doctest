## Preconditions
- The `ExtractBasePath` function is defined in `path_resolve`.
- Input is a pattern string like `"./foo/..."`.

## Steps
1. Call `path_resolve.ExtractBasePath(req.Input)`.
2. Return the result in `resp.StringResult`.

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
	result := path_resolve.ExtractBasePath(req.Input)
	return &Response{StringResult: result}, nil
}
```
