## Expected

- `CasesImportAssertPackage` returns true.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("detect failed: %v", err)
	}
	if !resp.Detected {
		t.Fatal("expected direct assert import to be detected")
	}
}
```