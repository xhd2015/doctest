## Expected

- `assertmod.ContentMD5()` returns non-empty hex string.
- Returned MD5 equals MD5 of `libdoc/assertmod/assert.go` on disk.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("assertmod md5 check failed: %v", err)
	}
	if resp.ContentMD5 == "" {
		t.Fatal("ContentMD5 returned empty string")
	}
	if resp.ContentMD5 != resp.FileMD5 {
		t.Fatalf("ContentMD5 mismatch: accessor=%s file=%s", resp.ContentMD5, resp.FileMD5)
	}
}
```