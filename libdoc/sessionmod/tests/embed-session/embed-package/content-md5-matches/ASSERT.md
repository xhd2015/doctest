## Expected

- `sessionmod.ContentMD5()` is non-empty hex.
- Matches MD5 of the embedded source file on disk.
- `Content()` length is greater than zero.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("sessionmod md5 check failed: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("sessionmod md5: %v", resp.Err)
	}
	if resp.ContentMD5 == "" {
		t.Fatal("ContentMD5 returned empty string")
	}
	if resp.ContentMD5 != resp.FileMD5 {
		t.Fatalf("ContentMD5 mismatch: accessor=%s file=%s", resp.ContentMD5, resp.FileMD5)
	}
	if resp.ContentLen == 0 {
		t.Fatal("Content() returned empty bytes")
	}
}
```
