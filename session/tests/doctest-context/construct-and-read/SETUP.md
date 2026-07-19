# Scenario

**Feature**: construct `session.Doctest` with all three fields and read each back

```
# composite literal sets ROOT, CASE, SESSION_ID
Caller -> session.Doctest{
    DOCTEST_ROOT:       "<abs root>",
    DOCTEST_CASE:       "<abs case>",
    DOCTEST_SESSION_ID: "<session id>",
}

# each field returns the value that was set
Caller <- d.DOCTEST_ROOT
Caller <- d.DOCTEST_CASE
Caller <- d.DOCTEST_SESSION_ID
```

## Preconditions

- All three `Want*` values are non-empty, distinct absolute-looking paths / id.
- Mode is `construct`.

## Steps

1. Set `WantRoot`, `WantCase`, and `WantSessionID` to distinct fixture strings.
2. Run constructs `session.Doctest` with those three fields.
3. Assert each observed field equals the corresponding `Want*`.

## Context

- Values are plain strings; absolute paths are representative, not resolved on disk.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "construct"
	req.WantRoot = "/tmp/doctest-context-root"
	req.WantCase = "/tmp/doctest-context-root/construct-and-read"
	req.WantSessionID = "doctest-context-construct-sid"
	return nil
}
```
