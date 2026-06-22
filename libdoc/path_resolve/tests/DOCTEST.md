# path_resolve Package Tests

## Version
0.0.2


These doc-style tests specify the contract for the `path_resolve` package.
Each function (`IsDotDotDotPattern`, `ExtractBasePath`, `ResolveRoot`,
`FindDotDotDotDirs`) is tested through direct function calls, not CLI
integration.

Tests are organized by function, with leaves covering all inputs and edge cases.

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
)

type Request struct {
	Input		string
	BasePath	string
}
type Response struct {
	BoolResult	bool
	StringResult	string
	DirsResult	[]string
	RootResult	string
	RootOkResult	bool
	ErrResult	string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	switch runType {
	case "extract_base_path":
		result := path_resolve.ExtractBasePath(req.Input)
		return &Response{StringResult: result}, nil
	case "find_dot_dot_dot_dirs":
		dirs, err := path_resolve.FindDotDotDotDirs(req.BasePath)
		if err != nil {
			return &Response{ErrResult: err.Error()}, err
		}
		return &Response{DirsResult: dirs}, nil
	case "is_dot_dot_dot_pattern":
		result := path_resolve.IsDotDotDotPattern(req.Input)
		return &Response{BoolResult: result}, nil
	case "resolve_root":
		root, ok := path_resolve.ResolveRoot(req.Input)
		return &Response{RootResult: root, RootOkResult: ok}, nil
	default:
		t.Fatalf("unknown runType: %s", runType)
		return nil, nil
	}
}
```
