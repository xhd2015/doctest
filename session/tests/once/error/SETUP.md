# Scenario

**Feature**: Once persists fn errors and replays them

```
# fn fails
Caller -> Once -> fn error "boom"
session.Once -> write error file; remove incomplete value
# second call
Caller -> Once -> return error without re-running fn
```

## Preconditions

- Valid session id and key.
- Mode `error` makes fn return `errors.New("boom")`.

## Steps

1. Call Once twice with the same key.
2. Assert both errors and single fn invocation.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "error"
	req.CallTwice = true
	return nil
}
```
