## Preconditions
- The `ResolveRoot` function is defined in `path_resolve`.
- It walks up the directory tree to find the root containing DOCTEST.md (or falls back to SETUP.md).

## Steps
1. Call `path_resolve.ResolveRoot(req.Input)`.
2. Return the result in `resp.RootResult` and `resp.RootOkResult`.

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
	root, ok := path_resolve.ResolveRoot(req.Input)
	return &Response{RootResult: root, RootOkResult: ok}, nil
}
```
