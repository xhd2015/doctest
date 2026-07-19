# Scenario

**Feature**: first fn error is persisted; second Once returns error without success

```
# first
Once -> fn returns boom -> error file
# second
Once -> read error -> return boom; fn not re-invoked
```

## Preconditions

- Unique session id for isolation.
- Key `"fail"`.

## Steps

1. CallTwice with Mode error.
2. Assert both errors contain boom; FnCalls == 1.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.SessionID = "once-doctest-err-" + d.DOCTEST_SESSION_ID
	req.Key = "fail"
	req.Mode = "error"
	req.CallTwice = true
	return nil
}
```
