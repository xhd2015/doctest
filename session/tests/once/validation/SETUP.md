# Scenario

**Feature**: Once rejects invalid inputs before invoking fn

```
# validation gates
session.Once -> missing DOCTEST_SESSION_ID => error (fn not run)
session.Once -> empty key => error (fn not run)
```

## Preconditions

- Validation happens before mkdir / flock / fn.
- Errors are non-nil; returned value is empty/nil.

## Steps

1. Leaf sets either empty session id or empty key.
2. Run calls Once once; Assert checks error and zero fn calls.

## Context

- Sibling leaves are MECE on the failing factor (session id vs key).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "json-object"
	req.CallTwice = false
	return nil
}
```
