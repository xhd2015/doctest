## Expected

- Both embed script runs succeed.
- First and second run MD5 digests are identical.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("embed script failed: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("embed script error: %v", resp.Err)
	}
	if resp.SecondRunMD5 == "" {
		t.Fatal("expected second run MD5 to be recorded")
	}
	if resp.OutputMD5 != resp.SecondRunMD5 {
		t.Fatalf("non-deterministic output: first=%s second=%s", resp.OutputMD5, resp.SecondRunMD5)
	}
}
```