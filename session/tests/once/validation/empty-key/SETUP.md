# Scenario

**Feature**: empty Once key is rejected

```
# session id present, key empty
Caller -> session.Once(t, "", fn)
Caller <- error empty key (fn not invoked)
```

## Preconditions

- Session id is set to a unique value for this leaf.
- Key is `""`.

## Steps

1. Set SessionID and empty Key.
2. Call Once once.
3. Assert error and FnCalls == 0.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.SessionID = "once-doctest-empty-key-" + d.DOCTEST_SESSION_ID
	req.Key = ""
	req.Mode = "json-object"
	return nil
}
```
