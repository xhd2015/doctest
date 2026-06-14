## Preconditions
- No --gen-dir is specified; build uses a temp directory.

## Steps
1. Create a project with 1 leaf.
2. Run `doctest build <test-dir> -v`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Log("using auto gen dir (mapping-gen cache, no --gen-dir)")
	return nil
}
```
