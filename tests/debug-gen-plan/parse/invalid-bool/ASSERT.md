## Expected

- Parse fails.
- Error is a **bool validation** failure for gen-plan (not merely "unknown key").

## Errors

- ParseErr non-empty; mentions gen-plan and invalid/bool semantics.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ParseErr == "" {
		t.Fatal("expected parse error for gen-plan=maybe")
	}
	// Must not stop at "unknown key" once gen-plan is implemented — require
	// bool-style validation so Classic TDD stays RED until key+parseBool land.
	low := strings.ToLower(resp.ParseErr)
	if strings.Contains(low, "unknown key") {
		t.Fatalf("want bool validation for gen-plan=maybe after key lands; still unknown key: %s", resp.ParseErr)
	}
	if !strings.Contains(low, "gen-plan") {
		t.Fatalf("expected gen-plan in error, got: %s", resp.ParseErr)
	}
	if !strings.Contains(low, "bool") && !strings.Contains(low, "invalid") {
		t.Fatalf("expected invalid-bool style error, got: %s", resp.ParseErr)
	}
}
```
