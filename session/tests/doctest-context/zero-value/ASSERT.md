## Expected

- `Run` succeeds with no harness error.
- `resp.View.Root` is `""`.
- `resp.View.Case` is `""`.
- `resp.View.SessionID` is `""`.

## Side Effects

- None.

## Errors

- None expected from `Run`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.View.Root != "" {
		t.Fatalf("zero DOCTEST_ROOT = %q, want empty", resp.View.Root)
	}
	if resp.View.Case != "" {
		t.Fatalf("zero DOCTEST_CASE = %q, want empty", resp.View.Case)
	}
	if resp.View.SessionID != "" {
		t.Fatalf("zero DOCTEST_SESSION_ID = %q, want empty", resp.View.SessionID)
	}
}
```
