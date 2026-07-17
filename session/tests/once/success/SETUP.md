# Scenario

**Feature**: successful Once stores and reuses raw JSON under the session cache

```
# valid sid + key
Caller -> session.Once -> fn(cacheDir) -> json.RawMessage
session.Once -> write value file
Caller <- same raw JSON on subsequent Once(same key)
```

## Preconditions

- Session id and key are non-empty.
- Fn returns valid JSON bytes and nil error.

## Steps

1. Leaf chooses session id, key(s), payload, and whether to call twice.
2. Assert checks equality, call counts, layout, or unmarshal.

## Context

- Success siblings split by observable behavior (reuse, independence, layout,
  client unmarshaling), not by trivial payload variants.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if req.SessionID == "" {
		req.SessionID = "once-doctest-success-" + DOCTEST_SESSION_ID
	}
	req.Mode = "json-object"
	return nil
}
```
