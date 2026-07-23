---
label: heavy
---

## Expected
- The test compiles and passes (import alias is preserved in generated code).
- Run returns without error since default Run just returns a success response.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected zero exit, got %d", resp.ExitCode)
	}
}
```
