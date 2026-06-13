## Preconditions
- The `FindDotDotDotDirs` function is defined in `path_resolve`.
- It discovers directories containing DOCTEST.md from a given base path.

## Steps
1. Call `path_resolve.FindDotDotDotDirs(req.BasePath)`.
2. On success, return dirs in `resp.DirsResult`.
3. On error, store the error message in `resp.ErrResult`.

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
	dirs, err := path_resolve.FindDotDotDotDirs(req.BasePath)
	if err != nil {
		return &Response{ErrResult: err.Error()}, err
	}
	return &Response{DirsResult: dirs}, nil
}
```
