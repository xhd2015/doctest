## Expected

- `Run` succeeds with no harness error.
- `resp.View.Root` equals `req.WantRoot`.
- `resp.View.Case` equals `req.WantCase`.
- `resp.View.SessionID` equals `req.WantSessionID`.

## Side Effects

- None. Construction and field reads are pure value operations.

## Errors

- None expected from `Run`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.View.Root != req.WantRoot {
		t.Fatalf("DOCTEST_ROOT = %q, want %q", resp.View.Root, req.WantRoot)
	}
	if resp.View.Case != req.WantCase {
		t.Fatalf("DOCTEST_CASE = %q, want %q", resp.View.Case, req.WantCase)
	}
	if resp.View.SessionID != req.WantSessionID {
		t.Fatalf("DOCTEST_SESSION_ID = %q, want %q", resp.View.SessionID, req.WantSessionID)
	}
}
```
