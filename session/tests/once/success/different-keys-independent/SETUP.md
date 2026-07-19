# Scenario

**Feature**: different Once keys under the same session are independent

```
# key A
Caller -> Once("a") -> fn -> {"key":...} for A path
# key B
Caller -> Once("b") -> fn again -> distinct JSON
```

## Preconditions

- Same session id for both keys.
- Keys `"a"` and `"b"`.

## Steps

1. Run Once for key A with default JSON object mode.
2. Run Once for key B with a distinct payload via SecondKey path.
3. Assert both succeed, values differ, FnCalls == 2.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.SessionID = "once-doctest-keys-" + d.DOCTEST_SESSION_ID
	req.Key = "a"
	req.SecondKey = "b"
	req.Mode = "json-object"
	req.JSONPayload = `{"key":"A"}`
	req.CallTwice = false
	return nil
}
```
