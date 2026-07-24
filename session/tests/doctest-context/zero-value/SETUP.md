# Scenario

**Feature**: zero value of `session.Doctest` has empty strings for all fields

```
# zero value — no composite literal fields set
Caller -> var d session.Doctest

# all three fields are empty strings
Caller <- d.DOCTEST_ROOT == ""
Caller <- d.DOCTEST_CASE == ""
Caller <- d.DOCTEST_SESSION_ID == ""
```

## Preconditions

- Mode is `zero`.
- No `Want*` values are required; construction does not assign fields.

## Steps

1. Set `req.Mode` to `zero`.
2. Run declares `var d session.Doctest` and snapshots the three fields.
3. Assert each observed field is the empty string.

## Context

- Confirms Go zero-value semantics for the public type once it exists.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "zero"
	// Want* intentionally left empty; zero path ignores them.
	return nil
}
```
