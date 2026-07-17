# Scenario

**Feature**: missing DOCTEST_SESSION_ID causes Once to error without running fn

```
# no session id in process env
Caller -> session.Once
session.Once <- syscall.Getenv empty/missing
Caller <- error (fn not invoked)
```

## Preconditions

- `req.SessionID` is empty so Run clears `DOCTEST_SESSION_ID`.
- Key is non-empty so the session-id check is the failing gate.

## Steps

1. Set SessionID empty, Key `"k"`.
2. Call Once once.
3. Assert error and FnCalls == 0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SessionID = ""
	req.Key = "k"
	req.Mode = "json-object"
	return nil
}
```
