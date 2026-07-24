## Expected

- Result matches `^nogit_[0-9a-f]{12}$`.
- Digest is SHA-256 of the exact abs root string, first 12 hex chars.

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := sha256.Sum256([]byte(req.AbsRoot))
	want := fmt.Sprintf("nogit_%s", hex.EncodeToString(sum[:])[:12])
	if resp.ProjectID != want {
		t.Fatalf("project id = %q, want %q", resp.ProjectID, want)
	}
	re := regexp.MustCompile(`^nogit_[0-9a-f]{12}$`)
	if !re.MatchString(resp.ProjectID) {
		t.Fatalf("project id %q does not match nogit_<12 hex>", resp.ProjectID)
	}
}
```
