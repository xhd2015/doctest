# Scenario

**Feature**: first Once writes a JSON object; second same key reuses without re-running fn

```
# first call
Caller -> Once(key) -> fn once -> write value
Caller <- raw JSON

# second call same session+key
Caller -> Once(key) -> read value (fn not run)
Caller <- identical raw JSON
```

## Preconditions

- Unique session id so this leaf does not reuse another leaf's value file.
- Same key for both calls.

## Steps

1. Set JSON payload `{"n":1,"label":"once"}`.
2. Call Once twice with the same key.
3. Assert equal raw bytes and FnCalls == 1.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SessionID = "once-doctest-reuse-" + DOCTEST_SESSION_ID
	req.Key = "cli-binary"
	req.Mode = "json-object"
	req.JSONPayload = `{"n":1,"label":"once"}`
	req.CallTwice = true
	return nil
}
```
