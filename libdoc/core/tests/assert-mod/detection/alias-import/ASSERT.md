## Expected

- `CasesImportAssertPackage` returns true for aliased import (path match, not alias name).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("detect failed: %v", err)
	}
	if !resp.Detected {
		t.Fatal("expected aliased assert import to be detected")
	}
}
```