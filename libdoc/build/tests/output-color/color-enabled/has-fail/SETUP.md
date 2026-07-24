# Scenario

**Feature**: mixed pass/fail colors fail dot and summary metrics

```
# colored regions (non-verbose only)
fail dot -> red | pass dot -> plain | summary Pass -> green | summary Fail -> red
```

## Preconditions
- One passing and one failing package (`a_pass_0` then `z_fail_0`).

## Steps
1. Set `PassCount` to 1 and `FailCount` to 1.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PassCount = 1
	req.FailCount = 1
	return nil
}
```