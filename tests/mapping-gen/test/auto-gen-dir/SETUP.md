## Preconditions
- No --gen-dir is specified; tests are generated under the mapping-gen cache root.

## Steps
1. Create a project with 1 leaf.
2. Run `doctest test <test-dir>`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Log("using auto gen dir (mapping-gen cache, no --gen-dir)")
	return nil
}
```
