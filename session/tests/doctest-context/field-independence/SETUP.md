# Scenario

**Feature**: setting one `session.Doctest` field does not change the other two

```
# three independent instances
Caller -> session.Doctest{ DOCTEST_ROOT: onlyRoot }       # CASE, SESSION_ID stay ""
Caller -> session.Doctest{ DOCTEST_CASE: onlyCase }       # ROOT, SESSION_ID stay ""
Caller -> session.Doctest{ DOCTEST_SESSION_ID: onlySID }  # ROOT, CASE stay ""
```

## Preconditions

- Mode is `independence`.
- `WantRoot`, `WantCase`, and `WantSessionID` are non-empty and pairwise distinct
  so a cross-field leak would be obvious.

## Steps

1. Set three distinct fixture values for Root, Case, and SessionID.
2. Run builds three separate `session.Doctest` values, each setting only one field.
3. Assert: set field equals Want*; the other two fields on that instance are `""`.

## Context

- Independence is about value semantics of the struct fields, not concurrent access.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "independence"
	req.WantRoot = "/tmp/doctest-indep-root"
	req.WantCase = "/tmp/doctest-indep-case"
	req.WantSessionID = "doctest-context-indep-sid"
	return nil
}
```
