## Expected

- `Run` succeeds with no harness error.
- **Only-root instance:** `Root == WantRoot`, `Case == ""`, `SessionID == ""`.
- **Only-case instance:** `Case == WantCase`, `Root == ""`, `SessionID == ""`.
- **Only-session instance:** `SessionID == WantSessionID`, `Root == ""`, `Case == ""`.

## Side Effects

- None. Three independent value constructions; no shared mutation.

## Errors

- None expected from `Run`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}

	// Set only DOCTEST_ROOT
	if resp.OnlyRoot.Root != req.WantRoot {
		t.Fatalf("only-root: DOCTEST_ROOT = %q, want %q", resp.OnlyRoot.Root, req.WantRoot)
	}
	if resp.OnlyRoot.Case != "" {
		t.Fatalf("only-root: DOCTEST_CASE = %q, want empty", resp.OnlyRoot.Case)
	}
	if resp.OnlyRoot.SessionID != "" {
		t.Fatalf("only-root: DOCTEST_SESSION_ID = %q, want empty", resp.OnlyRoot.SessionID)
	}

	// Set only DOCTEST_CASE
	if resp.OnlyCase.Case != req.WantCase {
		t.Fatalf("only-case: DOCTEST_CASE = %q, want %q", resp.OnlyCase.Case, req.WantCase)
	}
	if resp.OnlyCase.Root != "" {
		t.Fatalf("only-case: DOCTEST_ROOT = %q, want empty", resp.OnlyCase.Root)
	}
	if resp.OnlyCase.SessionID != "" {
		t.Fatalf("only-case: DOCTEST_SESSION_ID = %q, want empty", resp.OnlyCase.SessionID)
	}

	// Set only DOCTEST_SESSION_ID
	if resp.OnlySession.SessionID != req.WantSessionID {
		t.Fatalf("only-session: DOCTEST_SESSION_ID = %q, want %q", resp.OnlySession.SessionID, req.WantSessionID)
	}
	if resp.OnlySession.Root != "" {
		t.Fatalf("only-session: DOCTEST_ROOT = %q, want empty", resp.OnlySession.Root)
	}
	if resp.OnlySession.Case != "" {
		t.Fatalf("only-session: DOCTEST_CASE = %q, want empty", resp.OnlySession.Case)
	}
}
```
