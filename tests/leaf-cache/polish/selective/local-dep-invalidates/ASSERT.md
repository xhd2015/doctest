## Expected

- No harness error.
- Keys are lowercase hex digests.
- After local package mutation: **Key != Key2** (leaf key invalidated).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Err != "" {
		t.Fatalf("resp.Err: %s", resp.Err)
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("expected hex keys; Key=%q Key2=%q", resp.Key, resp.Key2)
	}
	if resp.Key == resp.Key2 {
		t.Fatalf("local dep edit must change leaf key; got same key %s", resp.Key)
	}
}
```
