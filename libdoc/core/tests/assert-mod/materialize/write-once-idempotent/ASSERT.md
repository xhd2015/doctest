## Expected

- Both `MaterializeAssertModule` calls succeed.
- `assert.go` MD5 digest unchanged after second call.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("materialize failed: %v", err)
	}
	afterDigest := snapshotMD5(t, filepath.Join(resp.CacheDir, "assert.go"))
	if afterDigest != beforeDigest {
		t.Fatal("assert.go content changed on second MaterializeAssertModule call")
	}
}
```