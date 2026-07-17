# Scenario

**Feature**: client can json.Unmarshal Once result into a struct

```
# typical testbin-style payload
fn -> {"path":"/tmp/example-bin"}
Caller <- json.RawMessage
json.Unmarshal -> struct { Path string `json:"path"` }
```

## Preconditions

- Payload is a JSON object with a `path` field (mirrors future testbin Ensure).

## Steps

1. Call Once once with `{"path":"/tmp/doctest-session-once-bin"}`.
2. Assert unmarshals into a typed struct with that path.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SessionID = "once-doctest-unmarshal-" + DOCTEST_SESSION_ID
	req.Key = "go-binary"
	req.Mode = "json-object"
	req.JSONPayload = `{"path":"/tmp/doctest-session-once-bin"}`
	req.CallTwice = false
	return nil
}
```
