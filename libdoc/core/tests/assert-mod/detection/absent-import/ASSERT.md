## Expected

- `CasesImportAssertPackage` returns false.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("detect failed: %v", err)
	}
	if resp.Detected {
		t.Fatal("expected no assert import to return false")
	}
}
```