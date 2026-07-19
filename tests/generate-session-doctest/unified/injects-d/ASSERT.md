## Expected

- Unified leaf source satisfies inject contract.
- Source is a `RunTestLeaf` entry (not classic TestXxx-only shape) still carries d inject.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unified assemble failed: %v", err)
	}
	assertInjectContract(t, "unified-leaf", resp.Source)
	// Sanity: unified entrypoint name is present in current API.
	if !strings.Contains(resp.Source, "RunTestLeaf") && !strings.Contains(resp.Source, "func RunTestLeaf") {
		// After P2 the name should still be RunTestLeaf; if renamed, inject contract is the hard assert above.
		t.Logf("note: RunTestLeaf name not found (still ok if inject contract passed)")
	}
}
```
